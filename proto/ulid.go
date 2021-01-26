// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package proto

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

func NewID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).String()
}

func NewIDBytes() []byte {
	id, _ := ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).MarshalText()
	return id
}

func WriteID(dst []byte) ([]byte, error) {
	err := ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).MarshalTextTo(dst)
	return dst, err
}

var globalRand = rand.New(&lockedSource{
	src: rand.NewSource(time.Now().UnixNano()),
})

type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
