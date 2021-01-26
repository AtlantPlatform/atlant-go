// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package proto

import (
	"bytes"
	"fmt"

	capn "github.com/glycerine/go-capnproto"

	"github.com/AtlantPlatform/atlant-go/state"
)

type AnnouncePeekFunc func(key *state.Key, v *Announce) error

func AnnouncePeek(fn AnnouncePeekFunc) state.PeekFunc {
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
		vv := ReadRootAnnounce(multiBuffer.Segments[0])
		return fn(k, &vv)
	}
}

type AnnounceModifyFunc func(key *state.Key, v *Announce) (*Announce, error)

func AnnounceModify(fn AnnounceModifyFunc) state.ModifyFunc {
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
			return v, err
		}
		vv := ReadRootAnnounce(seg)
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

func UnpackAnnounce(data []byte) (Announce, error) {
	seg, err := capn.ReadFromPackedStream(bytes.NewReader(data), nil)
	if err != nil {
		return Announce{}, err
	}
	v := ReadRootAnnounce(seg)
	return v, nil
}

type EnvelopeBeatTickPeekFunc func(key *state.Key, v *EnvelopeBeatTick) error

func EnvelopeBeatTickPeek(fn EnvelopeBeatTickPeekFunc) state.PeekFunc {
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
		vv := ReadRootEnvelopeBeatTick(multiBuffer.Segments[0])
		return fn(k, &vv)
	}
}

type EnvelopeBeatTickModifyFunc func(key *state.Key, v *EnvelopeBeatTick) (*EnvelopeBeatTick, error)

func EnvelopeBeatTickModify(fn EnvelopeBeatTickModifyFunc) state.ModifyFunc {
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
		vv := ReadRootEnvelopeBeatTick(seg)
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

func UnpackEnvelopeBeatTick(data []byte) (EnvelopeBeatTick, error) {
	seg, err := capn.ReadFromPackedStream(bytes.NewReader(data), nil)
	if err != nil {
		return EnvelopeBeatTick{}, err
	}
	v := ReadRootEnvelopeBeatTick(seg)
	return v, nil
}

type EnvelopeBeatInfoPeekFunc func(key *state.Key, v *EnvelopeBeatInfo) error

func EnvelopeBeatInfoPeek(fn EnvelopeBeatInfoPeekFunc) state.PeekFunc {
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
		vv := ReadRootEnvelopeBeatInfo(multiBuffer.Segments[0])
		return fn(k, &vv)
	}
}

type EnvelopeBeatInfoModifyFunc func(key *state.Key, v *EnvelopeBeatInfo) (*EnvelopeBeatInfo, error)

func EnvelopeBeatInfoModify(fn EnvelopeBeatInfoModifyFunc) state.ModifyFunc {
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
		vv := ReadRootEnvelopeBeatInfo(seg)
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

func UnpackEnvelopeBeatInfo(data []byte) (EnvelopeBeatInfo, error) {
	seg, err := capn.ReadFromPackedStream(bytes.NewReader(data), nil)
	if err != nil {
		return EnvelopeBeatInfo{}, err
	}
	v := ReadRootEnvelopeBeatInfo(seg)
	return v, nil
}

type EnvelopeRecordUpdatePeekFunc func(key *state.Key, v *EnvelopeRecordUpdate) error

func EnvelopeRecordUpdatePeek(fn EnvelopeRecordUpdatePeekFunc) state.PeekFunc {
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
		vv := ReadRootEnvelopeRecordUpdate(multiBuffer.Segments[0])
		return fn(k, &vv)
	}
}

type EnvelopeRecordUpdateModifyFunc func(key *state.Key, v *EnvelopeRecordUpdate) (*EnvelopeRecordUpdate, error)

func EnvelopeRecordUpdateModify(fn EnvelopeRecordUpdateModifyFunc) state.ModifyFunc {
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
		vv := ReadRootEnvelopeRecordUpdate(seg)
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

func UnpackEnvelopeRecordUpdate(data []byte) (EnvelopeRecordUpdate, error) {
	seg, err := capn.ReadFromPackedStream(bytes.NewReader(data), nil)
	if err != nil {
		return EnvelopeRecordUpdate{}, err
	}
	v := ReadRootEnvelopeRecordUpdate(seg)
	return v, nil
}
