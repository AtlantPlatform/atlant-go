// Copyright 2017-21 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package state

import "time"

type storeOptions struct {
	SyncWrites bool
	GCRatio    float64
	GCInterval time.Duration
}

type storeOpt func(o *storeOptions)

func defaultStoreOptions() *storeOptions {
	return &storeOptions{
		SyncWrites: true,
		GCRatio:    0.3,
		GCInterval: 5 * time.Minute,
	}
}

func NoSyncOption() storeOpt {
	return func(o *storeOptions) {
		o.SyncWrites = false
	}
}

func GCRatioOption(v float64) storeOpt {
	return func(o *storeOptions) {
		if v >= 0.01 && v <= 0.99 {
			o.GCRatio = v
		}
	}
}

func GCIntervalOption(v time.Duration) storeOpt {
	return func(o *storeOptions) {
		if v >= time.Second {
			o.GCInterval = v
		}
	}
}
