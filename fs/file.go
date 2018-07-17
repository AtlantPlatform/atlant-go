// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	capn "github.com/glycerine/go-capnproto"

	"github.com/AtlantPlatform/atlant-go/proto"
	"github.com/AtlantPlatform/go-ipfs/go-ipfs-cmdkit/files"
)

var (
	ErrNotDirectory = errors.New("could not get next file, not a directory")
	ErrNotReader    = errors.New("this is a directory, cannot call read functions")
)

func NewObjectFile(meta proto.ObjectMeta, body io.ReadCloser) (files.File, error) {
	metaBuf := new(bytes.Buffer)
	if _, err := meta.Segment.WriteToPacked(metaBuf); err != nil {
		err = fmt.Errorf("failed to pack object meta: %v", err)
		return nil, err
	}
	f := &objectFile{
		name: "object",
		pos:  0,
		files: []files.File{
			&metaFile{
				name: "meta",
				buf:  metaBuf,
				size: int64(metaBuf.Len()),
			},
		},
	}
	if body != nil {
		f.files = append(f.files, &contentFile{
			name: "content",
			body: body,
			size: meta.Size(),
		})
	}
	return f, nil
}

func readObjectFileMeta(body io.Reader) (proto.ObjectMeta, error) {
	seg, err := capn.ReadFromPackedStream(body, nil)
	if err != nil {
		return proto.ObjectMeta{}, err
	}
	meta := proto.ReadRootObjectMeta(seg)
	return meta, nil
}

// objectFile implements the files.File interface from IPFS.
type objectFile struct {
	name  string
	pos   int
	files []files.File
}

func (f *objectFile) IsDirectory() bool {
	return true
}

func (f *objectFile) NextFile() (files.File, error) {
	if f.pos >= len(f.files) {
		return nil, io.EOF
	}
	file := f.files[f.pos]
	f.pos++
	return file, nil
}

func (f *objectFile) FileName() string {
	return f.name
}

func (f *objectFile) FullPath() string {
	return f.name
}

func (f *objectFile) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (f *objectFile) Close() error {
	return ErrNotReader
}

func (f *objectFile) Peek(n int) files.File {
	return f.files[n]
}

func (f *objectFile) Length() int {
	return len(f.files)
}

func (f *objectFile) Size() (int64, error) {
	var size int64
	for _, file := range f.files {
		sizeFile, ok := file.(files.SizeFile)
		if !ok {
			return 0, errors.New("could not get size of child file")
		}
		s, err := sizeFile.Size()
		if err != nil {
			return 0, err
		}
		size += s
	}
	return size, nil
}

// metaFile implements the files.File interface from IPFS.
type metaFile struct {
	name string
	buf  *bytes.Buffer
	size int64
}

func (f *metaFile) IsDirectory() bool {
	return false
}

func (f *metaFile) NextFile() (files.File, error) {
	return nil, ErrNotDirectory
}

func (f *metaFile) FileName() string {
	return f.name
}

func (f *metaFile) FullPath() string {
	return f.name
}

func (f *metaFile) AbsPath() string {
	return f.name
}

func (f *metaFile) Read(p []byte) (int, error) {
	return f.buf.Read(p)
}

func (f *metaFile) Close() error {
	f.buf = nil
	return nil
}

func (f *metaFile) Size() (int64, error) {
	return f.size, nil
}

// contentFile implements the files.File interface from IPFS.
type contentFile struct {
	name string
	body io.ReadCloser
	size int64
}

func (f *contentFile) IsDirectory() bool {
	return false
}

func (f *contentFile) NextFile() (files.File, error) {
	return nil, ErrNotDirectory
}

func (f *contentFile) FileName() string {
	return f.name
}

func (f *contentFile) FullPath() string {
	return f.name
}

func (f *contentFile) AbsPath() string {
	return f.name
}

func (f *contentFile) Read(p []byte) (int, error) {
	return f.body.Read(p)
}

func (f *contentFile) Close() error {
	return f.body.Close()
}

func (f *contentFile) Size() (int64, error) {
	if f.size < 0 {
		return 0, errors.New("unknown content size")
	}
	return f.size, nil
}
