// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package state

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

// IndexedStore is an interface to internal node state.
// The state includes file names, version history, in other words it's a collection of indexed metadata.
// The state must be synchronised across all utility Atlant nodes.
type IndexedStore interface {
	View(k *Key, fn PeekFunc) error
	Update(k *Key, fn ModifyFunc) error
	Delete(k *Key) error

	RangeKeys(b Bucket, fn KeyFunc) (*RangeOptions, error)
	RangePeek(b Bucket, fn PeekFunc) (*RangeOptions, error)
	RangeModify(b Bucket, fn ModifyFunc) (*RangeOptions, error)

	Close() error
}

type BucketID uint16

type Bucket struct {
	ID           BucketID
	RangeOptions RangeOptions
}

func (b BucketID) Bytes() []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf[:2], uint16(b))
	return buf
}

func (b Bucket) WithRangeOptions(opt *RangeOptions) Bucket {
	return Bucket{
		ID:           b.ID,
		RangeOptions: *opt,
	}
}

var (
	BucketRecords   BucketID = 0x10
	BucketBeatTicks BucketID = 0x11
	BucketBeatInfos BucketID = 0x12
)

var NoKey = Bucket{}.NewKey(nil)

func (b Bucket) NewKey(key []byte) *Key {
	k := &Key{
		Bucket: b,
	}
	copy(k.Key[:], key)
	return k
}

type RangeOptions struct {
	Prefetch int
	Offset   []byte
	Limit    int
}

func NewBucket(id BucketID, opts ...*RangeOptions) Bucket {
	b := Bucket{ID: id}
	if len(opts) > 0 {
		if opts[0] != nil {
			b.RangeOptions = *opts[0]
		}
	}
	return b
}

type Key struct {
	Bucket Bucket
	Key    [26]byte
	TTL    time.Duration
}

func NewKey(bucket BucketID, key []byte) *Key {
	k := &Key{
		Bucket: Bucket{
			ID: bucket,
		},
	}
	copy(k.Key[:], key)
	return k
}

func (k *Key) WithinBucket(b Bucket) *Key {
	return &Key{
		Bucket: b,
		Key:    k.Key,
	}
}

func (k *Key) Bytes() []byte {
	buf := make([]byte, 2+26)
	binary.BigEndian.PutUint16(buf[:2], uint16(k.Bucket.ID))
	copy(buf[2:], k.Key[:])
	return buf
}

func (k *Key) Unmarshal(buf []byte) *Key {
	k.Bucket.ID = BucketID(binary.BigEndian.Uint16(buf[:2]))
	copy(k.Key[:], buf)
	return k
}

func (k *Key) String() string {
	return hex.EncodeToString(k.Bucket.ID.Bytes()) + string(k.Key[:])
}

var OffsetStart = []byte("")

var (
	ErrRangeStop = errors.New("range stop")
	ErrScanStop  = ErrRangeStop
)

var ErrNoUpdate = errors.New("no update")

type KeyFunc func(k *Key) error
type PeekFunc func(k *Key, v []byte) error
type ModifyFunc func(k *Key, v []byte) ([]byte, error)

func NewIndexedStoreBadger(prefix string, opts ...storeOpt) (IndexedStore, error) {
	return newBadgerStore(prefix, opts...)
}
