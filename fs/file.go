// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package fs

import (
	"bytes"
	"fmt"
	"io"

	capn "github.com/glycerine/go-capnproto"

	"github.com/AtlantPlatform/atlant-go/proto"
	files "github.com/ipfs/go-ipfs-files"
)

// NewObjectDir creates IPFS folder with "meta" and (optionally) "content" files inside
func NewObjectDir(meta proto.ObjectMeta, body io.Reader) (files.Directory, error) {
	metaBuf := new(bytes.Buffer)
	if _, err := meta.Segment.WriteToPacked(metaBuf); err != nil {
		err = fmt.Errorf("failed to pack object meta: %v", err)
		return nil, err
	}
	mapFiles := map[string]files.Node{
		"meta": files.NewBytesFile(metaBuf.Bytes()),
	}
	if body != nil {
		mapFiles["content"] = files.NewReaderFile(body)
	}
	return files.NewMapDirectory(mapFiles), nil
}

func readObjectFileMeta(body io.Reader) (proto.ObjectMeta, error) {
	seg, err := capn.ReadFromPackedStream(body, nil)
	if err != nil {
		return proto.ObjectMeta{}, err
	}
	meta := proto.ReadRootObjectMeta(seg)
	return meta, nil
}
