// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package fs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	capn "github.com/glycerine/go-capnproto"

	"github.com/AtlantPlatform/atlant-go/proto"
)

// Object is object in atlant node
type Object struct {
	ObjectRef

	Meta *proto.ObjectMeta
	Body io.ReadCloser
}

// PlanetaryFileStore interface
type PlanetaryFileStore interface {
	NodeID() string
	SignData(peerID string, data []byte) ([]byte, error)
	VerifyNode(code string) (string, error)

	PubSub() (PlanetaryPubSub, error)
	Listener() PlanetaryListener
	Client() PlanetaryClient

	PinObject(ref ObjectRef) error
	UnpinObject(ref ObjectRef) error
	PinNewest(ref ObjectRef, depth int) error

	PutObject(ctx context.Context, ref ObjectRef, userMeta []byte, body io.ReadCloser) (*ObjectRef, error)
	DeleteObject(ctx context.Context, ref ObjectRef) (*ObjectRef, error)
	GetObject(ctx context.Context, ref ObjectRef) (*Object, error)
	HeadObject(ctx context.Context, ref ObjectRef) (*ObjectRef, error)
	ListObjects(ctx context.Context, ref ObjectRef) ([]ObjectRef, error)

	DiskStats() (*DiskStats, error)
	BandwidthStats() *BandwidthStats
	RepoStats() *RepoStats
	BitswapStats() *BitswapStats

	Close() error
}

// NewPlanetaryFileStore creates new IPFS filestore
func NewPlanetaryFileStore(prefix string, opts ...IpfsOpt) (PlanetaryFileStore, error) {
	s, err := newIpfsStore(prefix, false, opts...)
	if err != nil {
		err = fmt.Errorf("failed to open IPFS store node: %v", err)
		return nil, err
	}
	return s, nil
}

// InitPlanetaryFileStore initializes new IPFS filestore
func InitPlanetaryFileStore(prefix string, opts ...IpfsOpt) (PlanetaryFileStore, error) {
	s, err := newIpfsStore(prefix, true, opts...)
	if err != nil {
		err = fmt.Errorf("failed to init IPFS store node: %v", err)
		return nil, err
	}
	return s, nil
}

// NewPrivateKey returns multiline buffer with swarm key
func NewPrivateKey() ([]byte, error) {
	r, err := newIpfsPrivateKey()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

// ObjectRef structure to keep meta information
type ObjectRef struct {
	ID   string
	Path string
	Size int64

	Version         string
	VersionPrevious string
	VersionOffset   int

	muxOnce sync.Once
	metaMux *sync.RWMutex
	meta    *proto.ObjectMeta
}

// SetMeta sets ObjectMeta with proper lock
func (o *ObjectRef) SetMeta(m *proto.ObjectMeta) {
	o.muxOnce.Do(func() {
		o.metaMux = new(sync.RWMutex)
	})
	o.metaMux.Lock()
	o.meta = m
	o.metaMux.Unlock()
}

// Meta gets ObjectMeta with proper read lock
func (o *ObjectRef) Meta() *proto.ObjectMeta {
	o.muxOnce.Do(func() {
		o.metaMux = new(sync.RWMutex)
	})
	o.metaMux.RLock()
	m := o.meta
	o.metaMux.RUnlock()
	return m
}

// ToProto converts object reference to ObjectMeta
func (o *ObjectRef) ToProto() (proto.ObjectMeta, error) {
	meta := proto.AutoNewObjectMeta(capn.NewBuffer(nil))
	if o == nil {
		return meta, nil
	}
	meta.SetId(o.ID)
	meta.SetPath(o.Path)
	if o.Size == 0 {
		meta.SetSize(-1)
	} else {
		meta.SetSize(o.Size)
	}
	meta.SetCreatedAt(time.Now().UnixNano())
	meta.SetVersionPrevious(o.VersionPrevious)
	return meta, nil
}

// PrevVersion decreases version offset
func (o *ObjectRef) PrevVersion() *ObjectRef {
	o.VersionOffset--
	return o
}

// NextVersion increases version offset
func (o *ObjectRef) NextVersion() *ObjectRef {
	o.VersionOffset++
	return o
}

// BandwidthStats - IPFS Bandwidth stats descriptor
type BandwidthStats struct {
	TotalIn  int64   `json:"total_in"`
	TotalOut int64   `json:"total_out"`
	RateIn   float64 `json:"rate_in"`
	RateOut  float64 `json:"rate_out"`
}

// DiskStats - IPFS Disk stats descriptor
type DiskStats struct {
	BytesAll  uint64 `json:"bytes_all"`
	BytesUsed uint64 `json:"bytes_used"`
	BytesFree uint64 `json:"bytes_free"`
}

// RepoStats - IPFS Repo stats descriptor
type RepoStats struct {
	NumObjects uint64 `json:"num_objects"`
	RepoSize   uint64 `json:"repo_size"`
	RepoPath   string `json:"repo_path"`
	Version    string `json:"version"`
	StorageMax uint64 `json:"storage_max"`
}

// BitswapStats - IPFS Bitswap descriptor
type BitswapStats struct {
	ProvideBufLen   int      `json:"provide_buf_len"`
	WantlistLen     int      `json:"wantlist_len"`
	Peers           []string `json:"peers"`
	BlocksReceived  uint64   `json:"blocks_received"`
	DataReceived    uint64   `json:"data_receiver"`
	BlocksSent      uint64   `json:"blocks_sent"`
	DataSent        uint64   `json:"data_sent"`
	DupBlksReceived uint64   `json:"dup_blks_received"`
	DupDataReceived uint64   `json:"dup_data_received"`
}
