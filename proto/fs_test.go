// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package proto

import (
	"bytes"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	"github.com/stretchr/testify/require"
)

func TestObjectMetaCapn(t *testing.T) {
	require := require.New(t)

	seg := capn.NewBuffer(nil)
	meta := AutoNewObjectMeta(seg)
	meta.SetPath("/test/hello.txt")
	buf := new(bytes.Buffer)
	meta.Segment.WriteToPacked(buf)

	segIn, err := capn.ReadFromPackedStream(buf, nil)
	require.NoError(err)
	metaOut := ReadRootObjectMeta(segIn)
	require.Equal("/test/hello.txt", metaOut.Path())
}
