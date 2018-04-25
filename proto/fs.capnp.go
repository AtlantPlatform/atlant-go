package proto

// AUTO GENERATED - DO NOT EDIT

import (
	"bufio"
	"bytes"
	"encoding/json"
	C "github.com/glycerine/go-capnproto"
	"io"
)

type ObjectMeta C.Struct

func NewObjectMeta(s *C.Segment) ObjectMeta      { return ObjectMeta(s.NewStruct(24, 5)) }
func NewRootObjectMeta(s *C.Segment) ObjectMeta  { return ObjectMeta(s.NewRootStruct(24, 5)) }
func AutoNewObjectMeta(s *C.Segment) ObjectMeta  { return ObjectMeta(s.NewStructAR(24, 5)) }
func ReadRootObjectMeta(s *C.Segment) ObjectMeta { return ObjectMeta(s.Root(0).ToStruct()) }
func (s ObjectMeta) Id() string                  { return C.Struct(s).GetObject(0).ToText() }
func (s ObjectMeta) IdBytes() []byte             { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s ObjectMeta) SetId(v string)              { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s ObjectMeta) Path() string                { return C.Struct(s).GetObject(1).ToText() }
func (s ObjectMeta) PathBytes() []byte           { return C.Struct(s).GetObject(1).ToDataTrimLastByte() }
func (s ObjectMeta) SetPath(v string)            { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s ObjectMeta) CreatedAt() int64            { return int64(C.Struct(s).Get64(0)) }
func (s ObjectMeta) SetCreatedAt(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s ObjectMeta) Version() string             { return C.Struct(s).GetObject(2).ToText() }
func (s ObjectMeta) VersionBytes() []byte        { return C.Struct(s).GetObject(2).ToDataTrimLastByte() }
func (s ObjectMeta) SetVersion(v string)         { C.Struct(s).SetObject(2, s.Segment.NewText(v)) }
func (s ObjectMeta) VersionPrevious() string     { return C.Struct(s).GetObject(3).ToText() }
func (s ObjectMeta) VersionPreviousBytes() []byte {
	return C.Struct(s).GetObject(3).ToDataTrimLastByte()
}
func (s ObjectMeta) SetVersionPrevious(v string) { C.Struct(s).SetObject(3, s.Segment.NewText(v)) }
func (s ObjectMeta) IsDeleted() bool             { return C.Struct(s).Get1(64) }
func (s ObjectMeta) SetIsDeleted(v bool)         { C.Struct(s).Set1(64, v) }
func (s ObjectMeta) Size() int64                 { return int64(C.Struct(s).Get64(16)) }
func (s ObjectMeta) SetSize(v int64)             { C.Struct(s).Set64(16, uint64(v)) }
func (s ObjectMeta) UserMeta() string            { return C.Struct(s).GetObject(4).ToText() }
func (s ObjectMeta) UserMetaBytes() []byte       { return C.Struct(s).GetObject(4).ToDataTrimLastByte() }
func (s ObjectMeta) SetUserMeta(v string)        { C.Struct(s).SetObject(4, s.Segment.NewText(v)) }
func (s ObjectMeta) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"id\":")
	if err != nil {
		return err
	}
	{
		s := s.Id()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"path\":")
	if err != nil {
		return err
	}
	{
		s := s.Path()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"createdAt\":")
	if err != nil {
		return err
	}
	{
		s := s.CreatedAt()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"version\":")
	if err != nil {
		return err
	}
	{
		s := s.Version()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"versionPrevious\":")
	if err != nil {
		return err
	}
	{
		s := s.VersionPrevious()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"isDeleted\":")
	if err != nil {
		return err
	}
	{
		s := s.IsDeleted()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"size\":")
	if err != nil {
		return err
	}
	{
		s := s.Size()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"userMeta\":")
	if err != nil {
		return err
	}
	{
		s := s.UserMeta()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s ObjectMeta) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s ObjectMeta) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('(')
	if err != nil {
		return err
	}
	_, err = b.WriteString("id = ")
	if err != nil {
		return err
	}
	{
		s := s.Id()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("path = ")
	if err != nil {
		return err
	}
	{
		s := s.Path()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("createdAt = ")
	if err != nil {
		return err
	}
	{
		s := s.CreatedAt()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("version = ")
	if err != nil {
		return err
	}
	{
		s := s.Version()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("versionPrevious = ")
	if err != nil {
		return err
	}
	{
		s := s.VersionPrevious()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("isDeleted = ")
	if err != nil {
		return err
	}
	{
		s := s.IsDeleted()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("size = ")
	if err != nil {
		return err
	}
	{
		s := s.Size()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("userMeta = ")
	if err != nil {
		return err
	}
	{
		s := s.UserMeta()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(')')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s ObjectMeta) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type ObjectMeta_List C.PointerList

func NewObjectMetaList(s *C.Segment, sz int) ObjectMeta_List {
	return ObjectMeta_List(s.NewCompositeList(24, 5, sz))
}
func (s ObjectMeta_List) Len() int            { return C.PointerList(s).Len() }
func (s ObjectMeta_List) At(i int) ObjectMeta { return ObjectMeta(C.PointerList(s).At(i).ToStruct()) }
func (s ObjectMeta_List) ToArray() []ObjectMeta {
	n := s.Len()
	a := make([]ObjectMeta, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s ObjectMeta_List) Set(i int, item ObjectMeta) { C.PointerList(s).Set(i, C.Object(item)) }
