// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package rs

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/y"
	capn "github.com/glycerine/go-capnproto"
	"github.com/oklog/ulid"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/authcenter"
	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/logging"
	"github.com/AtlantPlatform/atlant-go/proto"
	"github.com/AtlantPlatform/atlant-go/state"
)

// Record contains record descriptor
type Record struct {
	proto.Record

	Object fs.ObjectRef
	Body   io.ReadCloser
}

// RecordCRUD interface for Create/Read/Update/Delete record
type RecordCRUD interface {
	CreateRecord(ctx context.Context, path string, body io.ReadCloser, opts ...CreateOptions) (*Record, error)
	ReadRecord(ctx context.Context, path string, opts ...ReadOptions) (*Record, error)
	UpdateRecord(ctx context.Context, path string, body io.ReadCloser, opts ...UpdateOptions) (*Record, error)
	DeleteRecord(ctx context.Context, path string) (*Record, error)
}

// CreateOptions user meta
type CreateOptions struct {
	UserMeta []byte
	Size     int64
}

// UpdateOptions structure to contain user meta and size
type UpdateOptions struct {
	UserMeta []byte
	Size     int64
}

// ReadOptions structure - version
type ReadOptions struct {
	Version   string
	NoContent bool
}

// RecordWalkFunc handler to walk through path
type RecordWalkFunc func(path string, r *Record) error

// PlanetaryRecordStore interface to handle Record storage, see recordStore for implementation
type PlanetaryRecordStore interface {
	RecordCRUD

	ExportRecords(ctx context.Context, wr io.Writer) error
	WalkRecords(ctx context.Context, root string, fn RecordWalkFunc) error

	Sync(timeout time.Duration) error
	IsReady() bool
	WaitInbound(timeout time.Duration)
	WaitOutbound(timeout time.Duration)
	ReceiveEventAnnounce(event *EventAnnounce)
	EmitEventAnnounce(event *EventAnnounce)
	SendBeats(ctx context.Context, tickDur, infoDur time.Duration, ethAddr string)
	CommitBeatReports(ctx context.Context, dur time.Duration)

	BadgerStats() *BadgerStats
	Close() error
}

func GC(fileStore fs.PlanetaryFileStore, stateStore state.IndexedStore, maxVersions int) error {
	b := state.NewBucket(state.BucketRecords, &state.RangeOptions{
		Prefetch: 100,
	})
	_, err := stateStore.RangeModify(b, proto.RecordModify(func(k *state.Key, v *proto.Record) (*proto.Record, error) {
		if v.Previous().Len() <= maxVersions {
			return nil, state.ErrNoUpdate
		}
		rec := proto.AutoNewRecord(capn.NewBuffer(nil))
		rec.SetId(v.Id())
		rec.SetPath(v.Path())
		rec.SetCreatedAt(v.CreatedAt())
		ver := proto.AutoNewRecordVersion(capn.NewBuffer(nil))
		ver.SetAnnounce(v.Current().Announce())
		ver.SetVersion(v.Current().Version())
		rec.SetCurrent(ver)
		newPrevious, removed := proto.CapRecordVersions(v.Previous(), maxVersions)
		rec.SetPrevious(newPrevious)
		for _, versionID := range removed {
			if err := fileStore.UnpinObject(fs.ObjectRef{
				Version: versionID,
			}); err != nil {
				log.Debugln("failed to unpin during GC:", versionID, err)
			}
		}
		return &rec, nil
	}))
	return err
}

func NewPlanetaryRecordStore(nodeID string, fileStore fs.PlanetaryFileStore, stateStore state.IndexedStore) (PlanetaryRecordStore, error) {
	outboundAnnounces := make(chan *EventAnnounce, 1024)
	inboundAnnounces := make(chan *EventAnnounce, 1024)
	r := &recordStore{
		nodeID:   nodeID,
		stateMux: new(sync.RWMutex),

		fs: fileStore,
		ss: stateStore,

		outboundWg:        new(sync.WaitGroup),
		outboundPump:      pumpEventAnnounces(outboundAnnounces),
		outboundAnnounces: outboundAnnounces,

		inboundWg:        new(sync.WaitGroup),
		inboundPump:      pumpEventAnnounces(inboundAnnounces),
		inboundAnnounces: inboundAnnounces,
	}
	r.processInbound(4, 10*time.Minute)
	r.processOutbound(4, 10*time.Minute)

	sub, err := r.fs.PubSub()
	if err != nil {
		log.Warningf("failed to connect to pubsub: %v", err)
		return r, nil
	}
	pubsubConfig, err := sub.Config()
	if err != nil {
		log.Warningf("failed to find pubsub config: %v", err)
	}
	if !pubsubConfig.StrictSignatureVerification {
		log.Fatalln("Pubsub Config: StrictSignatureVerification is disabled. Please check fs/config")
	}

	topics := []string{
		EventRecordUpdate.String(),
		EventBeatInfo.String(),
		EventBeatTick.String(),
	}
	if err := sub.Subscribe(func(m *fs.Message) error {
		if m.From == r.nodeID {
			return nil
		} else if len(m.TopicIDs) == 0 {
			return nil
		}
		event := &EventAnnounce{
			Type: EventFromTopic(m.TopicIDs[0]),
		}
		switch event.Type {
		case EventUnknown:
			return nil
		case EventRecordUpdate:
			if !isPublishAllowed(m.From) {
				log.WithField("from", m.From).Debugln("Ignoring EventRecordUpdate, unauthorized node")
				return nil
			}
			log.WithFields(log.Fields{
				"from": m.From,
				"type": event.Type.String(),
			}).Debugln("EventRecordUpdate Received")

			seg, err := capn.ReadFromPackedStream(bytes.NewReader(m.Data), nil)
			if err != nil {
				log.WithFields(log.Fields{
					"from": m.From,
					"type": event.Type.String(),
				}).Warningln("Failed to decode record announce data:", err)
				return nil
			}
			event.Announce = proto.ReadRootAnnounce(seg)
			r.ReceiveEventAnnounce(event)
		case EventBeatTick, EventBeatInfo:
			seg, err := capn.ReadFromPackedStream(bytes.NewReader(m.Data), nil)
			if err != nil {
				log.WithFields(log.Fields{
					"from": m.From,
					"type": event.Type.String(),
				}).Warningln("Failed to decode beat announce data:", err)
				return nil
			}
			event.Announce = proto.ReadRootAnnounce(seg)
			r.ReceiveEventAnnounce(event)
		default:
			log.WithFields(log.Fields{
				"from": m.From,
				"type": event.Type.String(),
			}).Warningf("Event not handled")
			return nil
		}
		return nil
	}, topics...); err != nil {
		log.Warningln(err)
	}

	return r, nil
}

type recordStore struct {
	nodeID   string
	stateMux *sync.RWMutex
	state    storeState

	fs fs.PlanetaryFileStore
	ss state.IndexedStore

	outboundWg          *sync.WaitGroup
	outboundPump        chan *EventAnnounce
	outboundAnnounces   chan *EventAnnounce
	outboundWorkCounter uint64

	inboundWg          *sync.WaitGroup
	inboundPump        chan *EventAnnounce
	inboundAnnounces   chan *EventAnnounce
	inboundWorkCounter uint64
}

type storeState int

const (
	storeInactiveState storeState = 0
	storeSyncState     storeState = 1
	storeActiveState   storeState = 2
)

func (r *recordStore) Close() error {
	r.inboundPump <- &EventAnnounce{
		Type: EventStopAnnounce,
	}
	r.outboundPump <- &EventAnnounce{
		Type: EventStopAnnounce,
	}
	return nil
}

// ErrNotSynced to be thrown when not synced
var ErrNotSynced = errors.New("not synced")

func (r *recordStore) Sync(timeout time.Duration) error {
	var syncCandidates []string
	entries := authcenter.Default.Entries()
	for _, e := range entries {
		log.WithField("key", e.Key).Debug("Sync Candidate")
		if e.Key == r.nodeID {
			continue
		} else if e.HasPermissions(authcenter.RecordSyncPermission) {
			syncCandidates = append(syncCandidates, e.Key)
		}
	}
	if len(syncCandidates) == 0 {
		log.Warningln("no sync candidates found")
		r.state = storeActiveState
		return nil
	}
	log.Debugf("found %d sync candidates", len(syncCandidates))

	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()
	alive := r.aliveNodes(ctx, syncCandidates)
	if len(alive) == 0 {
		for i := 0; i < 3; i++ {
			log.Debugln("retrying to find alive candidates in 5s")
			time.Sleep(5 * time.Second)
			if alive = r.aliveNodes(ctx, syncCandidates); len(alive) > 0 {
				break
			}
		}
		if len(alive) == 0 {
			log.Warningln("no alive sync candidates found")
			r.state = storeActiveState
			return nil
		}
		log.Debugln("found alive sync candidates:", len(alive))
	} else {
		log.Debugln("found alive sync candidates:", len(alive))
	}
	if len(alive) > 2 {
		alive = alive[:2]
	}
	rC := make(chan *proto.Record, 100)
	go r.collectRecords(ctx, alive, rC)
	log.Infoln("sync has started with timeout", timeout)
	if err := r.startSync(ctx, rC); err != nil {
		err = fmt.Errorf("failed to sync store: %v", err)
		return err
	}
	return nil
}

func (r *recordStore) startSync(ctx context.Context, rC <-chan *proto.Record) error {
	r.setState(storeSyncState)
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				r.setState(storeActiveState)
				return nil
			}
			r.setState(storeInactiveState)
			return ErrNotSynced
		case record, ok := <-rC:
			if !ok {
				log.Debugln("sync end")
				r.setState(storeActiveState)
				return nil
			} else if err := validateRecord(record); err != nil {
				vv, _ := record.MarshalJSON()
				log.Debugf("failed to validate record in sync: %v, record: %s", err, string(vv))
				continue
			} else if ownerID := record.Current().Announce().NodeID(); !isPublishAllowed(ownerID) {
				log.Debugf("publish not allowed for author of the announce in sync: %s", ownerID)
				continue
			}
			k := state.NewKey(state.BucketRecords, record.IdBytes())
			if err := r.ss.Update(k, proto.RecordModify(func(k *state.Key, v *proto.Record) (*proto.Record, error) {
				if v == nil {
					// if not exists, simply insert
					log.Debugf("new record imported: %s", record.Id())
					return record, nil
				}
				updNext, err := record.AnnounceEnvelope()
				if err != nil {
					log.Debugf("failed to decode record update envelope in sync: %v", err)
					return nil, state.ErrNoUpdate
				}
				updCurrent, err := v.AnnounceEnvelope()
				if err != nil {
					log.Debugf("failed to decode current record in store: %v", err)
					return nil, state.ErrNoUpdate
				}
				if updNext.Id() != updCurrent.Id() {
					log.Warningf("announce envelope record ID mismatch: %s (next) != %s (prev)", updNext.Id(), updCurrent.Id())
					return nil, state.ErrNoUpdate
				}
				if cmp := updNext.Compare(updCurrent); cmp > 0 {
					// overwrite with new record, since its envelope is newer
					log.Debugf("record imported, newer version: %s", record.Id())
					return record, nil
				} else if cmp == 0 {
					// current envelopes are the same, compare lists
					if record.Previous().Len() > v.Previous().Len() {
						// overwrite if longer
						log.Debugf("record imported, version chain longer: %s", record.Id())
						return record, nil
					}
				}
				return nil, state.ErrNoUpdate
			})); err != nil {
				return err
			}
		}
	}
}

func validateRecord(record *proto.Record) error {
	if record == nil {
		return errors.New("record is nil")
	}
	ann := record.Current().Announce()
	ok, err := fs.VerifyDataSignature(ann.NodeID(), ann.Signature(), ann.Envelope())
	if err != nil {
		return fmt.Errorf("error checking current version signature: %v", err)
	} else if !ok {
		// fmt.Println("1. Node ID=", ann.NodeID())
		// fmt.Println("2. Signature=", ann.Signature())
		// fmt.Println("3. Envelope=", hex.EncodeToString(ann.Envelope()))
		return errors.New("incorrect signature for current version announce")
	}
	list := record.Previous()
	for i := 0; i < list.Len(); i++ {
		ann = list.At(i).Announce()
		ok, err := fs.VerifyDataSignature(ann.NodeID(), ann.Signature(), ann.Envelope())
		if err != nil {
			return fmt.Errorf("error checking version(%d) signature: %v", i, err)
		} else if !ok {
			return fmt.Errorf("incorrect signature for version(%d) announce", i)
		}
	}
	return nil
}

func (r *recordStore) outboundWork() {
	atomic.AddUint64(&r.outboundWorkCounter, 1)
}

func (r *recordStore) inboundWork() {
	atomic.AddUint64(&r.inboundWorkCounter, 1)
}

func (r *recordStore) processOutbound(workers int, emitTimeout time.Duration) {
	for i := 0; i < workers; i++ {
		r.outboundWg.Add(1)
		go func() {
			defer r.outboundWg.Done()

			for !r.IsReady() {
				time.Sleep(100 * time.Millisecond)
			}
			for ev := range r.outboundAnnounces {
				if err := r.emitEvent(ev, emitTimeout); err != nil {
					log.Warningln("error emitting event:", err)
				} else {
					r.outboundWork()
				}
			}
		}()
	}
}

func (r *recordStore) processInbound(workers int, timeout time.Duration) {
	for i := 0; i < workers; i++ {
		r.inboundWg.Add(1)
		go func() {
			defer r.inboundWg.Done()

			for !r.IsReady() {
				time.Sleep(100 * time.Millisecond)
			}
			for ev := range r.inboundAnnounces {
				if err := r.handleEvent(ev, timeout); err != nil {
					log.Warningln("error handling event:", err)
				} else {
					r.inboundWork()
				}
			}
		}()
	}
}

func (r *recordStore) SendBeats(ctx context.Context, tickDur, infoDur time.Duration, ethAddr string) {
	start := time.Now()
	session := ctx.Value("session_id").(string)
	tickTimer := time.NewTimer(tickDur)
	infoTimer := time.NewTimer(infoDur)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tickTimer.C:
			ann := r.newBeatTickAnnounce(session)
			r.EmitEventAnnounce(&EventAnnounce{
				Type:     EventBeatTick,
				Announce: *ann,
			})
			tickTimer.Reset(tickDur)
		case <-infoTimer.C:
			uptimeUnix := time.Since(start).Seconds()
			outboundWork := atomic.LoadUint64(&r.outboundWorkCounter)
			inboundWork := atomic.LoadUint64(&r.inboundWorkCounter)
			ann := r.newBeatInfoAnnounce(session, ethAddr, int64(uptimeUnix), outboundWork, inboundWork)
			r.EmitEventAnnounce(&EventAnnounce{
				Type:     EventBeatInfo,
				Announce: *ann,
			})
			infoTimer.Reset(infoDur)
		}
	}
}

// BeatReport should be the array of Sessions
type BeatReport struct {
	Sessions []*BeatSessionReport `json:"sessions"`
}

// BeatSessionReport should be session descriptor
type BeatSessionReport struct {
	SessionID    string `json:"session_id"`
	EthereumAddr string `json:"eth_addr"`
	Uptime       int    `json:"uptime_hours"`
	InboundWork  uint64 `json:"in_work"`
	OutboundWork uint64 `json:"out_work"`
}

func (r *recordStore) CommitBeatReports(ctx context.Context, dur time.Duration) {
	t := time.NewTimer(dur)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if !isPublishAllowed(r.nodeID) {
				t.Reset(dur)
				continue
			}
			reports := make(map[string]*BeatReport, 100)
			b := state.NewBucket(state.BucketBeatInfos)
			if _, err := r.ss.RangePeek(b,
				proto.EnvelopeBeatInfoPeek(func(k *state.Key, v *proto.EnvelopeBeatInfo) error {
					if v == nil {
						return nil
					}
					ethAddr := v.EthereumAddrBytes()
					if len(ethAddr) == 0 {
						return nil
					}
					report, ok := reports[string(ethAddr)]
					if !ok {
						report = &BeatReport{}
						reports[string(ethAddr)] = report
					}
					report.Sessions = append(report.Sessions, &BeatSessionReport{
						SessionID:    v.Session(),
						EthereumAddr: string(ethAddr),
						Uptime:       int(v.UptimeUnix() / 3600),
						InboundWork:  v.InboundWork(),
						OutboundWork: v.OutboundWork(),
					})
					return nil
				})); err != nil {
				log.Warningf("failed to count beat ticks: %v", err)
			}
			if !isPublishAllowed(r.nodeID) {
				t.Reset(dur)
				continue
			}

			buf := new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			for addr, report := range reports {
				if err := enc.Encode(report); err != nil {
					log.Errorf("failed to encode beat report: %v", err)
					return
				}
				exportPath := fmt.Sprintf("/beat_reports/%s.json", addr)
				_, err := r.CreateRecord(ctx, exportPath, ioutil.NopCloser(buf), CreateOptions{
					Size: int64(buf.Len()),
				})
				if err == ErrRecordExists {
					_, err = r.UpdateRecord(ctx, exportPath, ioutil.NopCloser(buf), UpdateOptions{
						Size: int64(buf.Len()),
					})
				}
				if err != nil {
					buf.Reset()
					log.Warningf("failed to write beat report to store: %v", err)
					time.Sleep(time.Second)
					continue
				}
				buf.Reset()
			}
			t.Reset(dur)
		}
	}
}

func (r *recordStore) emitEvent(ev *EventAnnounce, timeout time.Duration) error {
	// ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	// defer cancelFn()

	wg := new(sync.WaitGroup)
	defer wg.Wait()

	buf := new(bytes.Buffer)
	if _, err := ev.Announce.Segment.WriteToPacked(buf); err != nil {

		log.WithFields(log.Fields{
			"type":   ev.Type.String(),
			"nodeID": r.nodeID,
		}).Warningf("Failed to pack announce: %v", err)
		err = fmt.Errorf("Failed to pack announce: %v", err)
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		pub, err := r.fs.PubSub()
		if err != nil {
			log.WithFields(log.Fields{
				"type":   ev.Type.String(),
				"nodeID": r.nodeID,
			}).Warningf("Failed to use pubsub: %v", err)
			return
		}
		// topic := EventToTopic(r.nodeID, ev.Type)
		log.WithFields(log.Fields{
			"type":   ev.Type.String(),
			"nodeID": r.nodeID,
		}).Debugf("Emitting event to pubsub")

		if err = pub.Publish(ev.Type.String(), buf.Bytes()); err != nil {
			log.WithFields(log.Fields{
				"type":   ev.Type.String(),
				"nodeID": r.nodeID,
			}).Warningf("Pubsub publish failed: %v", err)
		}
	}()
	return nil
}

func isPublishAllowed(nodeID string) bool {
	return authcenter.Default.HasPermissions(nodeID, authcenter.RecordWritePermission)
}

var (
	defaultBeatTickTTL = 4 * time.Hour
	defaultBeatInfoTTL = 31 * 24 * time.Hour
)

func (r *recordStore) handleEvent(ev *EventAnnounce, timeout time.Duration) error {
	ownerID := ev.Announce.NodeID()
	fields := logging.WithFn(log.Fields{
		"OwnerID": ownerID,
	})
	if ownerID == r.nodeID {
		log.WithFields(fields).Debugln("skipping own event", ev.Type.String())
		return nil
	}
	validate := func(ev *EventAnnounce) bool {
		data := ev.Announce.Envelope()
		ok, err := fs.VerifyDataSignature(ownerID, ev.Announce.Signature(), data)
		if err != nil {
			log.WithFields(logging.WithMore(fields, log.Fields{
				"Signature": ev.Announce.Signature(),
				"DataLen":   len(data),
			})).Warningf("wrong signature: %v", err)
			return false
		} else if !ok {
			log.WithFields(logging.WithMore(fields, log.Fields{
				"Signature": ev.Announce.Signature(),
				"DataLen":   len(data),
			})).Warningf("record update signature not matching content")
			return false
		}
		return true
	}
	switch ev.Type {
	case EventRecordUpdate:
		if !isPublishAllowed(ownerID) {
			log.WithFields(fields).Warningf("skipping record update event from an unauthorized source")
			return nil
		} else if !validate(ev) {
			log.WithFields(fields).Warningf("skipping invalid record update event")
			return nil
		}
		update, err := proto.UnpackEnvelopeRecordUpdate(ev.Announce.Envelope())
		if err != nil {
			log.WithFields(fields).Errorf("failed to unpack record update: %v", err)
			return nil
		}
		ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
		ref, err := r.fs.HeadObject(ctx, fs.ObjectRef{
			Version: update.Version(),
		})
		cancelFn()
		updateFields := logging.WithMore(fields, log.Fields{
			"Version":     update.Version,
			"VersionPrev": update.VersionPrev,
		})
		if err == fs.ErrNotFound {
			log.WithFields(updateFields).Warningln("file not found on IPFS but announced")
			return nil
		} else if err != nil {
			log.WithFields(updateFields).Errorf("failed to retrieve object: %v", err)
			return nil
		}
		k := state.NewKey(state.BucketRecords, []byte(ref.ID))
		if err := r.ss.Update(k, proto.RecordModify(func(k *state.Key, v *proto.Record) (*proto.Record, error) {
			if v == nil {
				vv := proto.AutoNewRecord(capn.NewBuffer(nil))
				v = &vv
				v.SetId(ref.ID)
				v.SetPath(ref.Path)
				v.SetCreatedAt(ev.Announce.Timestamp())
				ver := proto.AutoNewRecordVersion(capn.NewBuffer(nil))
				ver.SetAnnounce(ev.Announce)
				ver.SetVersion(ref.Version)
				v.SetCurrent(ver)
				return v, nil
			}
			v.SetPrevious(proto.AppendRecordVersion(v.Previous(), v.Current()))
			ver := proto.AutoNewRecordVersion(capn.NewBuffer(nil))
			ver.SetAnnounce(ev.Announce)
			ver.SetVersion(ref.Version)
			v.SetCurrent(ver)
			return v, nil
		})); err != nil {
			log.Warningf("failed to update record: %v", err)
		}
		if err := r.fs.PinNewest(*ref, 3); err != nil {
			log.WithFields(updateFields).Errorf("failed to pin object: %v", err)
			return nil
		}
	case EventBeatTick:
		if !validate(ev) {
			log.WithFields(fields).Warningf("skipping invalid beat tick event")
			return nil
		}
		tick, err := proto.UnpackEnvelopeBeatTick(ev.Announce.Envelope())
		if err != nil {
			log.WithFields(fields).Errorf("failed to unpack beat tick: %v", err)
			return nil
		}
		k := state.NewKey(state.BucketBeatTicks, tick.IdBytes())
		k.TTL = defaultBeatTickTTL
		if err := r.ss.Update(k, proto.EnvelopeBeatTickModify(
			func(k *state.Key, v *proto.EnvelopeBeatTick) (*proto.EnvelopeBeatTick, error) {
				if v == nil {
					vv := proto.AutoNewEnvelopeBeatTick(capn.NewBuffer(nil))
					v = &vv
					v.SetId(tick.Id())
					v.SetSession(tick.Session())
					return v, nil
				}
				return nil, state.ErrNoUpdate
			})); err != nil {
			log.Warningf("failed to write tick: %v", err)
		}
	case EventBeatInfo:
		if !validate(ev) {
			log.WithFields(fields).Warningf("skipping invalid beat info event")
			return nil
		}
		info, err := proto.UnpackEnvelopeBeatInfo(ev.Announce.Envelope())
		if err != nil {
			log.WithFields(fields).Errorf("failed to unpack beat info: %v", err)
			return nil
		}
		u, err := ulid.Parse(info.Id())
		if err != nil {
			log.WithFields(fields).Errorf("failed to parse beat info timestamp: %v", err)
			return nil
		} else if l := len(info.EthereumAddrBytes()); l == 0 || l > 64 {
			log.WithFields(fields).Errorf("skipping beat with incorrect eth address length: %d", l)
			return nil
		}
		lowerBound := u.Time() - uint64(info.UptimeUnix()*1000)
		var ticks int
		b := state.NewBucket(state.BucketBeatTicks)
		if _, err := r.ss.RangePeek(b,
			proto.EnvelopeBeatTickPeek(func(k *state.Key, v *proto.EnvelopeBeatTick) error {
				if v == nil {
					return nil
				}
				u, err := ulid.Parse(v.Id())
				if err != nil {
					return nil
				} else if u.Time() < lowerBound {
					// ignore ticks before uptime started
					return nil
				}
				if bytes.Equal(v.SessionBytes(), info.SessionBytes()) {
					ticks++
				}
				return nil
			})); err != nil {
			log.Warningf("failed to count beat ticks: %v", err)
		}
		k := state.NewKey(state.BucketBeatInfos, info.SessionBytes())
		k.TTL = defaultBeatInfoTTL
		if err := r.ss.Update(k, proto.EnvelopeBeatInfoModify(
			func(k *state.Key, v *proto.EnvelopeBeatInfo) (*proto.EnvelopeBeatInfo, error) {
				if v == nil {
					if ticks == 0 {
						// no prior ticks
						return nil, state.ErrNoUpdate
					}
					vv := proto.AutoNewEnvelopeBeatInfo(capn.NewBuffer(nil))
					v = &vv
					v.SetId(info.Id())
					v.SetSession(info.Session())
					v.SetEthereumAddr(info.EthereumAddr())
					v.SetUptimeUnix(info.UptimeUnix())
					v.SetOutboundWork(info.OutboundWork())
					v.SetInboundWork(info.InboundWork())
					return v, nil
				} else if info.UptimeUnix() > v.UptimeUnix() {
					if info.EthereumAddr() != v.EthereumAddr() {
						// same session, different addr? go away
						return nil, state.ErrNoUpdate
					} else if ticks < 3 {
						return nil, state.ErrNoUpdate
					}
					v.SetUptimeUnix(info.UptimeUnix())
					v.SetOutboundWork(info.OutboundWork())
					v.SetInboundWork(info.InboundWork())
					return v, nil
				}
				return nil, state.ErrNoUpdate
			})); err != nil {
			log.Warningf("failed to write beat info: %v", err)
		}
	default:
		log.Warningln("skipping unknown event:", ev.Type.String())
	}
	return nil
}

func (r *recordStore) IsReady() bool {
	r.stateMux.RLock()
	ready := r.state == storeActiveState
	r.stateMux.RUnlock()
	return ready
}

func (r *recordStore) setState(state storeState) {
	r.stateMux.Lock()
	r.state = state
	r.stateMux.Unlock()
}

func (r *recordStore) WaitOutbound(timeout time.Duration) {
	waitWG(r.outboundWg, timeout)
}

func (r *recordStore) WaitInbound(timeout time.Duration) {
	waitWG(r.inboundWg, timeout)
}

func waitWG(wg *sync.WaitGroup, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		select {
		case <-done:
		default:
			close(done)
		}
	}()
	select {
	case <-time.Tick(timeout):
	case <-done:
	}
}

// ReceiveEventAnnounce never blocks. Internal workers will eventually handle the received events.
func (r *recordStore) ReceiveEventAnnounce(event *EventAnnounce) {
	if event.Type == EventStopAnnounce {
		return
	}
	r.inboundPump <- event
}

// EmitEventAnnounce never blocks. Internal workers will eventually handle the events to emit.
func (r *recordStore) EmitEventAnnounce(event *EventAnnounce) {
	if event.Type == EventStopAnnounce {
		return
	}
	r.outboundPump <- event
}

var (
	// ErrNotAuthorized to be thrown when node is not authorized
	ErrNotAuthorized = errors.New("node is not authorized to create records")
	// ErrRecordExists to be thrown when record already exists
	ErrRecordExists = errors.New("record exists")
	// ErrRecordNotFound to be thrown when record was not found
	ErrRecordNotFound = errors.New("record not found")
)

func (r *recordStore) CreateRecord(ctx context.Context, path string, body io.ReadCloser, opts ...CreateOptions) (*Record, error) {
	if !isPublishAllowed(r.nodeID) {
		return nil, ErrNotAuthorized
	}
	defer r.inboundWork()
	id, err := r.findRecordID(ctx, path, "")
	if len(id) > 0 {
		return nil, ErrRecordExists
	} else if err != ErrRecordNotFound {
		return nil, err
	}
	id = proto.NewID()
	k := state.NewKey(state.BucketRecords, []byte(id))
	var size int64
	var userMeta []byte
	if len(opts) > 0 {
		size = opts[0].Size
		userMeta = opts[0].UserMeta
	}

	var ann *proto.Announce
	rec := &Record{}
	if err := r.ss.Update(k, proto.RecordModify(func(k *state.Key, v *proto.Record) (*proto.Record, error) {
		if v != nil {
			return v, ErrRecordExists
		}
		ref, err := r.fs.PutObject(ctx, fs.ObjectRef{
			ID:   id,
			Path: path,
			Size: size,
		}, userMeta, body)

		if err != nil {
			log.WithFields(log.Fields{
				"id":       id,
				"path":     path,
				"size":     size,
				"userMeta": string(userMeta),
			}).Errorf("IPFS error of PutObject: (CreateRecord) %v", err)
			return nil, err
		}
		log.WithFields(log.Fields{
			"id":       id,
			"path":     path,
			"size":     size,
			"userMeta": string(userMeta),
		}).Info("IPFS PutObject on CreateRecord was successfull")

		ann = r.newRecordUpdateAnnounce(id, ref.Version, "")
		rec.Record = proto.AutoNewRecord(capn.NewBuffer(nil))
		rec.Record.SetId(ref.ID)
		rec.Record.SetPath(ref.Path)
		rec.Record.SetCreatedAt(ann.Timestamp())
		ver := proto.AutoNewRecordVersion(capn.NewBuffer(nil))
		ver.SetAnnounce(*ann)
		ver.SetVersion(ref.Version)
		rec.Record.SetCurrent(ver)
		rec.Object = *ref
		return &rec.Record, nil
	})); err != nil {
		log.Errorf("failed to update record: %v", err)
		return nil, err
	} else if ann != nil {
		r.EmitEventAnnounce(&EventAnnounce{
			Type:     EventRecordUpdate,
			Announce: *ann,
		})
	} else {
		log.Errorln("record updated but the announce is empty")
	}
	return rec, nil
}

func (r *recordStore) findRecordID(ctx context.Context, path, version string) (string, error) {
	if len(version) > 0 {
		if ref, err := r.fs.HeadObject(ctx, fs.ObjectRef{
			Version: version,
		}); err == nil && len(ref.ID) > 0 {
			return ref.ID, err
		}
	}
	if u, err := ulid.Parse(path); err == nil && u.Time() > 0 {
		// path parsed as a valid ULID
		return path, nil
	}
	b := state.NewBucket(state.BucketRecords)
	var id string
	_, err := r.ss.RangePeek(b, proto.RecordPeek(func(k *state.Key, v *proto.Record) error {
		if v.Path() == path {
			id = v.Id()
			return state.ErrRangeStop
		}
		return nil
	}))
	if err != nil {
		return "", err
	} else if len(id) == 0 {
		return "", ErrRecordNotFound
	}
	return id, nil
}

func (r *recordStore) UpdateRecord(ctx context.Context, path string, body io.ReadCloser, opts ...UpdateOptions) (*Record, error) {
	if !isPublishAllowed(r.nodeID) {
		return nil, ErrNotAuthorized
	}
	defer r.inboundWork()
	id, err := r.findRecordID(ctx, path, "")
	if err != nil {
		return nil, err
	}
	k := state.NewKey(state.BucketRecords, []byte(id))
	var size int64
	var userMeta []byte
	if len(opts) > 0 {
		size = opts[0].Size
		userMeta = opts[0].UserMeta
	}

	var ann *proto.Announce
	rec := &Record{}
	if err := r.ss.Update(k, proto.RecordModify(func(k *state.Key, v *proto.Record) (*proto.Record, error) {
		if v == nil {
			return nil, ErrRecordNotFound
		}
		ref, err := r.fs.PutObject(ctx, fs.ObjectRef{
			ID:              v.Id(),
			Path:            path,
			VersionPrevious: v.Current().Version(),
			Size:            size,
		}, userMeta, body)
		if err != nil {
			log.WithFields(log.Fields{
				"id":       id,
				"path":     path,
				"size":     size,
				"userMeta": string(userMeta),
			}).Errorf("IPFS error of PutObject (UpdateRecord): %v", err)
			return nil, err
		}
		log.WithFields(log.Fields{
			"id":       id,
			"path":     path,
			"size":     size,
			"userMeta": string(userMeta),
		}).Info("IPFS PutObject on UpdateRecord was successfull")

		ann = r.newRecordUpdateAnnounce(id, ref.Version, v.Current().Version())
		v.SetPrevious(proto.AppendRecordVersion(v.Previous(), v.Current()))
		ver := proto.AutoNewRecordVersion(capn.NewBuffer(nil))
		ver.SetAnnounce(*ann)
		ver.SetVersion(ref.Version)
		v.SetCurrent(ver)
		rec.Record = *v
		rec.Object = *ref
		return v, nil
	})); err != nil {
		log.Errorf("failed to update record: %v", err)
		return nil, err
	} else if ann != nil {
		r.EmitEventAnnounce(&EventAnnounce{
			Type:     EventRecordUpdate,
			Announce: *ann,
		})
	} else {
		log.Errorln("record updated but the announce is empty")
	}
	return rec, nil
}

func (r *recordStore) DeleteRecord(ctx context.Context, path string) (*Record, error) {
	if !isPublishAllowed(r.nodeID) {
		return nil, ErrNotAuthorized
	}
	defer r.inboundWork()
	id, err := r.findRecordID(ctx, path, "")
	if err != nil {
		return nil, err
	}
	k := state.NewKey(state.BucketRecords, []byte(id))

	var ann *proto.Announce
	rec := &Record{}
	if err := r.ss.Update(k, proto.RecordModify(func(k *state.Key, v *proto.Record) (*proto.Record, error) {
		if v == nil {
			return nil, ErrRecordNotFound
		}
		if ref, err := r.fs.HeadObject(ctx, fs.ObjectRef{
			Version: v.Current().Version(),
		}); err == fs.ErrNotFound {
			return nil, ErrRecordNotFound
		} else if err != nil {
			return nil, err
		} else if ref.Meta().IsDeleted() {
			rec.Record = *v
			rec.Object = *ref
			return nil, nil
		}
		ref, err := r.fs.DeleteObject(ctx, fs.ObjectRef{
			ID:              v.Id(),
			Path:            v.Path(),
			VersionPrevious: v.Current().Version(),
		})
		if err != nil {
			return nil, err
		}
		ann = r.newRecordUpdateAnnounce(id, ref.Version, v.Current().Version())
		v.SetPrevious(proto.AppendRecordVersion(v.Previous(), v.Current()))
		ver := proto.AutoNewRecordVersion(capn.NewBuffer(nil))
		ver.SetAnnounce(*ann)
		ver.SetVersion(ref.Version)
		v.SetCurrent(ver)
		rec.Record = *v
		rec.Object = *ref
		return v, nil
	})); err != nil {
		log.Errorf("failed to update record: %v", err)
		return nil, err
	}
	if ann != nil {
		r.EmitEventAnnounce(&EventAnnounce{
			Type:     EventRecordUpdate,
			Announce: *ann,
		})
	}
	return rec, nil
}

func (r *recordStore) newBeatTickAnnounce(session string) *proto.Announce {
	e := proto.AutoNewEnvelopeBeatTick(capn.NewBuffer(nil))
	e.SetId(proto.NewID())
	e.SetSession(session)
	buf := new(bytes.Buffer)
	if _, err := e.Segment.WriteToPacked(buf); err != nil {
		panic(fmt.Sprintf("failed to pack data: %v", err))
	}
	sig, err := r.fs.SignData(r.nodeID, buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("failed to use FS signer: %+v", err.Error()))
	}
	a := proto.AutoNewAnnounce(capn.NewBuffer(nil))
	a.SetId(proto.NewID())
	a.SetType(proto.ANNOUNCETYPE_BEATTICK)
	a.SetEnvelope(buf.Bytes())
	a.SetSignature(hex.EncodeToString(sig))
	a.SetTimestamp(time.Now().UnixNano())
	a.SetNodeID(r.nodeID)
	return &a
}

func (r *recordStore) newBeatInfoAnnounce(session string, ethAddr string, uptimeUnix int64, announcesN, requestsN uint64) *proto.Announce {
	e := proto.AutoNewEnvelopeBeatInfo(capn.NewBuffer(nil))
	e.SetId(proto.NewID())
	e.SetSession(session)
	e.SetEthereumAddr(ethAddr)
	e.SetUptimeUnix(uptimeUnix)
	e.SetOutboundWork(announcesN)
	e.SetInboundWork(requestsN)
	buf := new(bytes.Buffer)
	if _, err := e.Segment.WriteToPacked(buf); err != nil {
		panic(fmt.Sprintf("failed to pack data: %v", err))
	}
	sig, err := r.fs.SignData(r.nodeID, buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("failed to use FS signer: %v", err))
	}
	a := proto.AutoNewAnnounce(capn.NewBuffer(nil))
	a.SetId(proto.NewID())
	a.SetType(proto.ANNOUNCETYPE_BEATINFO)
	a.SetEnvelope(buf.Bytes())
	a.SetSignature(hex.EncodeToString(sig))
	a.SetTimestamp(time.Now().UnixNano())
	a.SetNodeID(r.nodeID)
	return &a
}

func (r *recordStore) newRecordUpdateAnnounce(id, ver, verPrev string) *proto.Announce {
	e := proto.AutoNewEnvelopeRecordUpdate(capn.NewBuffer(nil))
	e.SetId(id)
	e.SetVersion(ver)
	e.SetVersionPrev(verPrev)
	buf := new(bytes.Buffer)
	if _, err := e.Segment.WriteToPacked(buf); err != nil {
		panic(fmt.Sprintf("failed to pack data: %v", err))
	}
	sig, err := r.fs.SignData(r.nodeID, buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("failed to use FS signer: %v", err))
	}
	a := proto.AutoNewAnnounce(capn.NewBuffer(nil))
	a.SetId(proto.NewID())
	a.SetType(proto.ANNOUNCETYPE_RECORDUPDATE)
	a.SetEnvelope(buf.Bytes())
	a.SetSignature(hex.EncodeToString(sig))
	a.SetTimestamp(time.Now().UnixNano())
	a.SetNodeID(r.nodeID)
	return &a
}

func (r *recordStore) ReadRecord(ctx context.Context, path string, opts ...ReadOptions) (*Record, error) {
	var version string
	var noContent bool
	if len(opts) > 0 {
		version = opts[0].Version
		noContent = opts[0].NoContent
	}
	defer r.inboundWork()
	id, err := r.findRecordID(ctx, path, version)
	if err != nil {
		return nil, err
	}
	k := state.NewKey(state.BucketRecords, []byte(id))
	rec := &Record{}
	var reqVersion string
	if err = r.ss.View(k, proto.RecordPeek(func(k *state.Key, v *proto.Record) error {
		if v == nil {
			return ErrRecordNotFound
		}
		// careful: a deep copy might be required
		rec.Record = *v
		reqVersion = v.Current().Version()
		return nil
	})); err != nil && err != ErrRecordNotFound {
		log.Warningln(err)
	}
	if len(version) > 0 {
		reqVersion = version
	}
	if len(reqVersion) == 0 {
		return nil, ErrRecordNotFound
	}
	if noContent {
		ref, err := r.fs.HeadObject(ctx, fs.ObjectRef{
			Version: reqVersion,
		})
		if err == fs.ErrNotFound {
			return nil, ErrRecordNotFound
		} else if err != nil {
			return nil, err
		} else if ref.Meta().IsDeleted() {
			rec.Object = *ref
			return rec, ErrRecordNotFound
		}
		rec.Object = *ref
	} else {
		obj, err := r.fs.GetObject(ctx, fs.ObjectRef{
			Version: reqVersion,
		})
		if err == fs.ErrNotFound {
			return nil, ErrRecordNotFound
		} else if err != nil {
			return nil, err
		} else if obj.Meta.IsDeleted() {
			rec.Object = obj.ObjectRef
			return rec, ErrRecordNotFound
		}
		rec.Object = obj.ObjectRef
		rec.Body = obj.Body
	}
	return rec, nil
}

// ErrWalkStop should be thrown to stop records traversing
var ErrWalkStop = errors.New("walk stop")

func (r *recordStore) WalkRecords(ctx context.Context, root string, fn RecordWalkFunc) error {
	defer r.inboundWork()
	b := state.NewBucket(state.BucketRecords, &state.RangeOptions{
		Offset: []byte(root),
	})
	_, err := r.ss.RangePeek(b, proto.RecordPeek(func(k *state.Key, v *proto.Record) error {
		if err := fn(v.Path(), &Record{
			Record: *v,
		}); err == ErrWalkStop {
			return state.ErrRangeStop
		} else if err != nil {
			return err
		}
		return nil
	}))
	return err
}

func (r *recordStore) ExportRecords(ctx context.Context, wr io.Writer) error {
	defer r.inboundWork()
	b := state.NewBucket(state.BucketRecords, &state.RangeOptions{
		Prefetch: 100,
	})
	_, err := r.ss.RangePeek(b, func(k *state.Key, v []byte) error {
		_, err := io.Copy(wr, bytes.NewReader(v))
		if err == io.EOF {
			return state.ErrRangeStop
		} else if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *recordStore) BadgerStats() *BadgerStats {
	return &BadgerStats{
		NumReads:        y.NumReads.Value(),
		NumWrites:       y.NumWrites.Value(),
		NumBytesRead:    y.NumBytesRead.Value(),
		NumBytesWritten: y.NumBytesWritten.Value(),
		NumLSMGets:      json.RawMessage(y.NumLSMGets.String()),
		NumLSMBloomHits: json.RawMessage(y.NumLSMBloomHits.String()),
		NumGets:         y.NumGets.Value(),
		NumPuts:         y.NumPuts.Value(),
		NumBlockedPuts:  y.NumBlockedPuts.Value(),
		NumMemtableGets: y.NumMemtableGets.Value(),
		LSMSize:         json.RawMessage(y.LSMSize.String()),
		VlogSize:        json.RawMessage(y.VlogSize.String()),
		PendingWrites:   json.RawMessage(y.PendingWrites.String()),
	}
}

// BadgerStats contains stats of badger
type BadgerStats struct {
	NumReads        int64           `json:"badger_disk_reads_total"`
	NumWrites       int64           `json:"badger_disk_writes_total"`
	NumBytesRead    int64           `json:"badger_read_bytes"`
	NumBytesWritten int64           `json:"badger_written_bytes"`
	NumLSMGets      json.RawMessage `json:"badger_lsm_level_gets_total"`
	NumLSMBloomHits json.RawMessage `json:"badger_lsm_bloom_hits_total"`
	NumGets         int64           `json:"badger_gets_total"`
	NumPuts         int64           `json:"badger_puts_total"`
	NumBlockedPuts  int64           `json:"badger_blocked_puts_total"`
	NumMemtableGets int64           `json:"badger_memtable_gets_total"`
	LSMSize         json.RawMessage `json:"badger_lsm_size_bytes"`
	VlogSize        json.RawMessage `json:"badger_vlog_size_bytes"`
	PendingWrites   json.RawMessage `json:"badger_pending_writes_total"`
}
