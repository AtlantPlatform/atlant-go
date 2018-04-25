package proto

import (
	"bytes"
	"fmt"

	capn "github.com/glycerine/go-capnproto"
	"github.com/oklog/ulid"

	"github.com/AtlantPlatform/atlant-go/state"
)

func (r *Record) AnnounceEnvelope() (*EnvelopeRecordUpdate, error) {
	seg, err := capn.ReadFromPackedStream(bytes.NewReader(r.Current().Announce().Envelope()), nil)
	if err != nil {
		return nil, err
	}
	vv := ReadRootEnvelopeRecordUpdate(seg)
	return &vv, nil
}

func (e *EnvelopeRecordUpdate) Compare(e2 *EnvelopeRecordUpdate) int {
	if e.Id() == e2.Id() {
		return 0
	} else if len(e.Id()) == 0 && len(e2.Id()) > 0 {
		return -1
	} else if len(e.Id()) > 0 && len(e2.Id()) == 0 {
		return 1
	}
	id, err := ulid.Parse(e.Id())
	id2, err2 := ulid.Parse(e.Id())
	if err != nil && err2 != nil {
		return 0
	} else if err != nil && err2 == nil {
		return -1
	} else if err == nil && err2 != nil {
		return 1
	}
	return id.Compare(id2)
}

type RecordPeekFunc func(key *state.Key, v *Record) error

func RecordPeek(fn RecordPeekFunc) state.PeekFunc {
	return func(k *state.Key, v []byte) error {
		if v == nil {
			return fn(k, nil)
		}
		multiBuffer := capn.NewSingleSegmentMultiBuffer()
		read, err := capn.ReadFromMemoryZeroCopyNoAlloc(v, multiBuffer)
		if err != nil {
			return err
		} else if read != int64(len(v)) {
			panic(fmt.Sprintf("wrong read: %d != %d", read, len(v)))
		}
		vv := ReadRootRecord(multiBuffer.Segments[0])
		return fn(k, &vv)
	}
}

type RecordModifyFunc func(key *state.Key, v *Record) (*Record, error)

func RecordModify(fn RecordModifyFunc) state.ModifyFunc {
	return func(k *state.Key, v []byte) ([]byte, error) {
		if v == nil {
			ret, err := fn(k, nil)
			if err != nil || ret == nil {
				return nil, err
			}
			buf := new(bytes.Buffer)
			if _, err := ret.Segment.WriteTo(buf); err != nil {
				return v, err
			}
			return buf.Bytes(), nil
		}
		seg, err := capn.ReadFromStream(bytes.NewReader(v), nil)
		if err != nil {
			return nil, err
		}
		vv := ReadRootRecord(seg)
		ret, err := fn(k, &vv)
		if err != nil || ret == nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		if _, err := ret.Segment.WriteTo(buf); err != nil {
			return v, err
		}
		return buf.Bytes(), nil
	}
}

func AppendRecordVersion(list RecordVersion_List, ver RecordVersion) RecordVersion_List {
	newList := NewRecordVersionList(capn.NewBuffer(nil), list.Len()+1)
	prevArr := list.ToArray()
	for i := range prevArr {
		newList.Set(i, prevArr[i])
	}
	newList.Set(list.Len(), ver)
	return newList
}
