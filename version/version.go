// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package version

// Version components
const (
	Maj = "1"
	Min = "0"
	Fix = "1-rc1"
)

var (
	// Version is the current version of Atlant Node
	// Must be a string because scripts like dist.sh read this file.
	Version = "1.0.1-rc1"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
