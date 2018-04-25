package proto

// AUTO GENERATED - DO NOT EDIT

import (
	"bufio"
	"bytes"
	"encoding/json"
	C "github.com/glycerine/go-capnproto"
	"io"
)

type Record C.Struct

func NewRecord(s *C.Segment) Record               { return Record(s.NewStruct(8, 4)) }
func NewRootRecord(s *C.Segment) Record           { return Record(s.NewRootStruct(8, 4)) }
func AutoNewRecord(s *C.Segment) Record           { return Record(s.NewStructAR(8, 4)) }
func ReadRootRecord(s *C.Segment) Record          { return Record(s.Root(0).ToStruct()) }
func (s Record) Id() string                       { return C.Struct(s).GetObject(0).ToText() }
func (s Record) IdBytes() []byte                  { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s Record) SetId(v string)                   { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Record) Path() string                     { return C.Struct(s).GetObject(1).ToText() }
func (s Record) PathBytes() []byte                { return C.Struct(s).GetObject(1).ToDataTrimLastByte() }
func (s Record) SetPath(v string)                 { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s Record) CreatedAt() int64                 { return int64(C.Struct(s).Get64(0)) }
func (s Record) SetCreatedAt(v int64)             { C.Struct(s).Set64(0, uint64(v)) }
func (s Record) Current() RecordVersion           { return RecordVersion(C.Struct(s).GetObject(2).ToStruct()) }
func (s Record) SetCurrent(v RecordVersion)       { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Record) Previous() RecordVersion_List     { return RecordVersion_List(C.Struct(s).GetObject(3)) }
func (s Record) SetPrevious(v RecordVersion_List) { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Record) WriteJSON(w io.Writer) error {
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
	_, err = b.WriteString("\"current\":")
	if err != nil {
		return err
	}
	{
		s := s.Current()
		err = s.WriteJSON(b)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"previous\":")
	if err != nil {
		return err
	}
	{
		s := s.Previous()
		{
			err = b.WriteByte('[')
			if err != nil {
				return err
			}
			for i, s := range s.ToArray() {
				if i != 0 {
					_, err = b.WriteString(", ")
				}
				if err != nil {
					return err
				}
				err = s.WriteJSON(b)
				if err != nil {
					return err
				}
			}
			err = b.WriteByte(']')
		}
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
func (s Record) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s Record) WriteCapLit(w io.Writer) error {
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
	_, err = b.WriteString("current = ")
	if err != nil {
		return err
	}
	{
		s := s.Current()
		err = s.WriteCapLit(b)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("previous = ")
	if err != nil {
		return err
	}
	{
		s := s.Previous()
		{
			err = b.WriteByte('[')
			if err != nil {
				return err
			}
			for i, s := range s.ToArray() {
				if i != 0 {
					_, err = b.WriteString(", ")
				}
				if err != nil {
					return err
				}
				err = s.WriteCapLit(b)
				if err != nil {
					return err
				}
			}
			err = b.WriteByte(']')
		}
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
func (s Record) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type Record_List C.PointerList

func NewRecordList(s *C.Segment, sz int) Record_List { return Record_List(s.NewCompositeList(8, 4, sz)) }
func (s Record_List) Len() int                       { return C.PointerList(s).Len() }
func (s Record_List) At(i int) Record                { return Record(C.PointerList(s).At(i).ToStruct()) }
func (s Record_List) ToArray() []Record {
	n := s.Len()
	a := make([]Record, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s Record_List) Set(i int, item Record) { C.PointerList(s).Set(i, C.Object(item)) }

type RecordVersion C.Struct

func NewRecordVersion(s *C.Segment) RecordVersion      { return RecordVersion(s.NewStruct(0, 2)) }
func NewRootRecordVersion(s *C.Segment) RecordVersion  { return RecordVersion(s.NewRootStruct(0, 2)) }
func AutoNewRecordVersion(s *C.Segment) RecordVersion  { return RecordVersion(s.NewStructAR(0, 2)) }
func ReadRootRecordVersion(s *C.Segment) RecordVersion { return RecordVersion(s.Root(0).ToStruct()) }
func (s RecordVersion) Version() string                { return C.Struct(s).GetObject(0).ToText() }
func (s RecordVersion) VersionBytes() []byte           { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s RecordVersion) SetVersion(v string)            { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s RecordVersion) Announce() Announce             { return Announce(C.Struct(s).GetObject(1).ToStruct()) }
func (s RecordVersion) SetAnnounce(v Announce)         { C.Struct(s).SetObject(1, C.Object(v)) }
func (s RecordVersion) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
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
	_, err = b.WriteString("\"announce\":")
	if err != nil {
		return err
	}
	{
		s := s.Announce()
		err = s.WriteJSON(b)
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
func (s RecordVersion) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s RecordVersion) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('(')
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
	_, err = b.WriteString("announce = ")
	if err != nil {
		return err
	}
	{
		s := s.Announce()
		err = s.WriteCapLit(b)
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
func (s RecordVersion) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type RecordVersion_List C.PointerList

func NewRecordVersionList(s *C.Segment, sz int) RecordVersion_List {
	return RecordVersion_List(s.NewCompositeList(0, 2, sz))
}
func (s RecordVersion_List) Len() int { return C.PointerList(s).Len() }
func (s RecordVersion_List) At(i int) RecordVersion {
	return RecordVersion(C.PointerList(s).At(i).ToStruct())
}
func (s RecordVersion_List) ToArray() []RecordVersion {
	n := s.Len()
	a := make([]RecordVersion, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s RecordVersion_List) Set(i int, item RecordVersion) { C.PointerList(s).Set(i, C.Object(item)) }

type Announce C.Struct

func NewAnnounce(s *C.Segment) Announce      { return Announce(s.NewStruct(16, 4)) }
func NewRootAnnounce(s *C.Segment) Announce  { return Announce(s.NewRootStruct(16, 4)) }
func AutoNewAnnounce(s *C.Segment) Announce  { return Announce(s.NewStructAR(16, 4)) }
func ReadRootAnnounce(s *C.Segment) Announce { return Announce(s.Root(0).ToStruct()) }
func (s Announce) Id() string                { return C.Struct(s).GetObject(0).ToText() }
func (s Announce) IdBytes() []byte           { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s Announce) SetId(v string)            { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Announce) NodeID() string            { return C.Struct(s).GetObject(1).ToText() }
func (s Announce) NodeIDBytes() []byte       { return C.Struct(s).GetObject(1).ToDataTrimLastByte() }
func (s Announce) SetNodeID(v string)        { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s Announce) Signature() string         { return C.Struct(s).GetObject(2).ToText() }
func (s Announce) SignatureBytes() []byte    { return C.Struct(s).GetObject(2).ToDataTrimLastByte() }
func (s Announce) SetSignature(v string)     { C.Struct(s).SetObject(2, s.Segment.NewText(v)) }
func (s Announce) Timestamp() int64          { return int64(C.Struct(s).Get64(0)) }
func (s Announce) SetTimestamp(v int64)      { C.Struct(s).Set64(0, uint64(v)) }
func (s Announce) Type() AnnounceType        { return AnnounceType(C.Struct(s).Get16(8)) }
func (s Announce) SetType(v AnnounceType)    { C.Struct(s).Set16(8, uint16(v)) }
func (s Announce) Envelope() []byte          { return C.Struct(s).GetObject(3).ToData() }
func (s Announce) SetEnvelope(v []byte)      { C.Struct(s).SetObject(3, s.Segment.NewData(v)) }
func (s Announce) WriteJSON(w io.Writer) error {
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
	_, err = b.WriteString("\"nodeID\":")
	if err != nil {
		return err
	}
	{
		s := s.NodeID()
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
	_, err = b.WriteString("\"signature\":")
	if err != nil {
		return err
	}
	{
		s := s.Signature()
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
	_, err = b.WriteString("\"timestamp\":")
	if err != nil {
		return err
	}
	{
		s := s.Timestamp()
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
	_, err = b.WriteString("\"type\":")
	if err != nil {
		return err
	}
	{
		s := s.Type()
		err = s.WriteJSON(b)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"envelope\":")
	if err != nil {
		return err
	}
	{
		s := s.Envelope()
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
func (s Announce) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s Announce) WriteCapLit(w io.Writer) error {
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
	_, err = b.WriteString("nodeID = ")
	if err != nil {
		return err
	}
	{
		s := s.NodeID()
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
	_, err = b.WriteString("signature = ")
	if err != nil {
		return err
	}
	{
		s := s.Signature()
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
	_, err = b.WriteString("timestamp = ")
	if err != nil {
		return err
	}
	{
		s := s.Timestamp()
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
	_, err = b.WriteString("type = ")
	if err != nil {
		return err
	}
	{
		s := s.Type()
		err = s.WriteCapLit(b)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("envelope = ")
	if err != nil {
		return err
	}
	{
		s := s.Envelope()
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
func (s Announce) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type Announce_List C.PointerList

func NewAnnounceList(s *C.Segment, sz int) Announce_List {
	return Announce_List(s.NewCompositeList(16, 4, sz))
}
func (s Announce_List) Len() int          { return C.PointerList(s).Len() }
func (s Announce_List) At(i int) Announce { return Announce(C.PointerList(s).At(i).ToStruct()) }
func (s Announce_List) ToArray() []Announce {
	n := s.Len()
	a := make([]Announce, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s Announce_List) Set(i int, item Announce) { C.PointerList(s).Set(i, C.Object(item)) }

type AnnounceType uint16

const (
	ANNOUNCETYPE_UNKNOWN      AnnounceType = 0
	ANNOUNCETYPE_BEATTICK     AnnounceType = 1
	ANNOUNCETYPE_BEATINFO     AnnounceType = 2
	ANNOUNCETYPE_RECORDUPDATE AnnounceType = 3
)

func (c AnnounceType) String() string {
	switch c {
	case ANNOUNCETYPE_UNKNOWN:
		return "unknown"
	case ANNOUNCETYPE_BEATTICK:
		return "beatTick"
	case ANNOUNCETYPE_BEATINFO:
		return "beatInfo"
	case ANNOUNCETYPE_RECORDUPDATE:
		return "recordUpdate"
	default:
		return ""
	}
}

func AnnounceTypeFromString(c string) AnnounceType {
	switch c {
	case "unknown":
		return ANNOUNCETYPE_UNKNOWN
	case "beatTick":
		return ANNOUNCETYPE_BEATTICK
	case "beatInfo":
		return ANNOUNCETYPE_BEATINFO
	case "recordUpdate":
		return ANNOUNCETYPE_RECORDUPDATE
	default:
		return 0
	}
}

type AnnounceType_List C.PointerList

func NewAnnounceTypeList(s *C.Segment, sz int) AnnounceType_List {
	return AnnounceType_List(s.NewUInt16List(sz))
}
func (s AnnounceType_List) Len() int              { return C.UInt16List(s).Len() }
func (s AnnounceType_List) At(i int) AnnounceType { return AnnounceType(C.UInt16List(s).At(i)) }
func (s AnnounceType_List) ToArray() []AnnounceType {
	n := s.Len()
	a := make([]AnnounceType, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s AnnounceType_List) Set(i int, item AnnounceType) { C.UInt16List(s).Set(i, uint16(item)) }
func (s AnnounceType) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	buf, err = json.Marshal(s.String())
	if err != nil {
		return err
	}
	_, err = b.Write(buf)
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s AnnounceType) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s AnnounceType) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	_, err = b.WriteString(s.String())
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s AnnounceType) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type EnvelopeBeatTick C.Struct

func NewEnvelopeBeatTick(s *C.Segment) EnvelopeBeatTick { return EnvelopeBeatTick(s.NewStruct(0, 2)) }
func NewRootEnvelopeBeatTick(s *C.Segment) EnvelopeBeatTick {
	return EnvelopeBeatTick(s.NewRootStruct(0, 2))
}
func AutoNewEnvelopeBeatTick(s *C.Segment) EnvelopeBeatTick {
	return EnvelopeBeatTick(s.NewStructAR(0, 2))
}
func ReadRootEnvelopeBeatTick(s *C.Segment) EnvelopeBeatTick {
	return EnvelopeBeatTick(s.Root(0).ToStruct())
}
func (s EnvelopeBeatTick) Id() string           { return C.Struct(s).GetObject(0).ToText() }
func (s EnvelopeBeatTick) IdBytes() []byte      { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s EnvelopeBeatTick) SetId(v string)       { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s EnvelopeBeatTick) Session() string      { return C.Struct(s).GetObject(1).ToText() }
func (s EnvelopeBeatTick) SessionBytes() []byte { return C.Struct(s).GetObject(1).ToDataTrimLastByte() }
func (s EnvelopeBeatTick) SetSession(v string)  { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s EnvelopeBeatTick) WriteJSON(w io.Writer) error {
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
	_, err = b.WriteString("\"session\":")
	if err != nil {
		return err
	}
	{
		s := s.Session()
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
func (s EnvelopeBeatTick) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s EnvelopeBeatTick) WriteCapLit(w io.Writer) error {
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
	_, err = b.WriteString("session = ")
	if err != nil {
		return err
	}
	{
		s := s.Session()
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
func (s EnvelopeBeatTick) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type EnvelopeBeatTick_List C.PointerList

func NewEnvelopeBeatTickList(s *C.Segment, sz int) EnvelopeBeatTick_List {
	return EnvelopeBeatTick_List(s.NewCompositeList(0, 2, sz))
}
func (s EnvelopeBeatTick_List) Len() int { return C.PointerList(s).Len() }
func (s EnvelopeBeatTick_List) At(i int) EnvelopeBeatTick {
	return EnvelopeBeatTick(C.PointerList(s).At(i).ToStruct())
}
func (s EnvelopeBeatTick_List) ToArray() []EnvelopeBeatTick {
	n := s.Len()
	a := make([]EnvelopeBeatTick, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s EnvelopeBeatTick_List) Set(i int, item EnvelopeBeatTick) {
	C.PointerList(s).Set(i, C.Object(item))
}

type EnvelopeBeatInfo C.Struct

func NewEnvelopeBeatInfo(s *C.Segment) EnvelopeBeatInfo { return EnvelopeBeatInfo(s.NewStruct(24, 3)) }
func NewRootEnvelopeBeatInfo(s *C.Segment) EnvelopeBeatInfo {
	return EnvelopeBeatInfo(s.NewRootStruct(24, 3))
}
func AutoNewEnvelopeBeatInfo(s *C.Segment) EnvelopeBeatInfo {
	return EnvelopeBeatInfo(s.NewStructAR(24, 3))
}
func ReadRootEnvelopeBeatInfo(s *C.Segment) EnvelopeBeatInfo {
	return EnvelopeBeatInfo(s.Root(0).ToStruct())
}
func (s EnvelopeBeatInfo) Id() string           { return C.Struct(s).GetObject(0).ToText() }
func (s EnvelopeBeatInfo) IdBytes() []byte      { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s EnvelopeBeatInfo) SetId(v string)       { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s EnvelopeBeatInfo) Session() string      { return C.Struct(s).GetObject(1).ToText() }
func (s EnvelopeBeatInfo) SessionBytes() []byte { return C.Struct(s).GetObject(1).ToDataTrimLastByte() }
func (s EnvelopeBeatInfo) SetSession(v string)  { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s EnvelopeBeatInfo) EthereumAddr() string { return C.Struct(s).GetObject(2).ToText() }
func (s EnvelopeBeatInfo) EthereumAddrBytes() []byte {
	return C.Struct(s).GetObject(2).ToDataTrimLastByte()
}
func (s EnvelopeBeatInfo) SetEthereumAddr(v string) { C.Struct(s).SetObject(2, s.Segment.NewText(v)) }
func (s EnvelopeBeatInfo) UptimeUnix() int64        { return int64(C.Struct(s).Get64(0)) }
func (s EnvelopeBeatInfo) SetUptimeUnix(v int64)    { C.Struct(s).Set64(0, uint64(v)) }
func (s EnvelopeBeatInfo) InboundWork() uint64      { return C.Struct(s).Get64(8) }
func (s EnvelopeBeatInfo) SetInboundWork(v uint64)  { C.Struct(s).Set64(8, v) }
func (s EnvelopeBeatInfo) OutboundWork() uint64     { return C.Struct(s).Get64(16) }
func (s EnvelopeBeatInfo) SetOutboundWork(v uint64) { C.Struct(s).Set64(16, v) }
func (s EnvelopeBeatInfo) WriteJSON(w io.Writer) error {
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
	_, err = b.WriteString("\"session\":")
	if err != nil {
		return err
	}
	{
		s := s.Session()
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
	_, err = b.WriteString("\"ethereumAddr\":")
	if err != nil {
		return err
	}
	{
		s := s.EthereumAddr()
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
	_, err = b.WriteString("\"uptimeUnix\":")
	if err != nil {
		return err
	}
	{
		s := s.UptimeUnix()
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
	_, err = b.WriteString("\"inboundWork\":")
	if err != nil {
		return err
	}
	{
		s := s.InboundWork()
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
	_, err = b.WriteString("\"outboundWork\":")
	if err != nil {
		return err
	}
	{
		s := s.OutboundWork()
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
func (s EnvelopeBeatInfo) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s EnvelopeBeatInfo) WriteCapLit(w io.Writer) error {
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
	_, err = b.WriteString("session = ")
	if err != nil {
		return err
	}
	{
		s := s.Session()
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
	_, err = b.WriteString("ethereumAddr = ")
	if err != nil {
		return err
	}
	{
		s := s.EthereumAddr()
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
	_, err = b.WriteString("uptimeUnix = ")
	if err != nil {
		return err
	}
	{
		s := s.UptimeUnix()
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
	_, err = b.WriteString("inboundWork = ")
	if err != nil {
		return err
	}
	{
		s := s.InboundWork()
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
	_, err = b.WriteString("outboundWork = ")
	if err != nil {
		return err
	}
	{
		s := s.OutboundWork()
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
func (s EnvelopeBeatInfo) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type EnvelopeBeatInfo_List C.PointerList

func NewEnvelopeBeatInfoList(s *C.Segment, sz int) EnvelopeBeatInfo_List {
	return EnvelopeBeatInfo_List(s.NewCompositeList(24, 3, sz))
}
func (s EnvelopeBeatInfo_List) Len() int { return C.PointerList(s).Len() }
func (s EnvelopeBeatInfo_List) At(i int) EnvelopeBeatInfo {
	return EnvelopeBeatInfo(C.PointerList(s).At(i).ToStruct())
}
func (s EnvelopeBeatInfo_List) ToArray() []EnvelopeBeatInfo {
	n := s.Len()
	a := make([]EnvelopeBeatInfo, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s EnvelopeBeatInfo_List) Set(i int, item EnvelopeBeatInfo) {
	C.PointerList(s).Set(i, C.Object(item))
}

type EnvelopeRecordUpdate C.Struct

func NewEnvelopeRecordUpdate(s *C.Segment) EnvelopeRecordUpdate {
	return EnvelopeRecordUpdate(s.NewStruct(0, 3))
}
func NewRootEnvelopeRecordUpdate(s *C.Segment) EnvelopeRecordUpdate {
	return EnvelopeRecordUpdate(s.NewRootStruct(0, 3))
}
func AutoNewEnvelopeRecordUpdate(s *C.Segment) EnvelopeRecordUpdate {
	return EnvelopeRecordUpdate(s.NewStructAR(0, 3))
}
func ReadRootEnvelopeRecordUpdate(s *C.Segment) EnvelopeRecordUpdate {
	return EnvelopeRecordUpdate(s.Root(0).ToStruct())
}
func (s EnvelopeRecordUpdate) Id() string      { return C.Struct(s).GetObject(0).ToText() }
func (s EnvelopeRecordUpdate) IdBytes() []byte { return C.Struct(s).GetObject(0).ToDataTrimLastByte() }
func (s EnvelopeRecordUpdate) SetId(v string)  { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s EnvelopeRecordUpdate) Version() string { return C.Struct(s).GetObject(1).ToText() }
func (s EnvelopeRecordUpdate) VersionBytes() []byte {
	return C.Struct(s).GetObject(1).ToDataTrimLastByte()
}
func (s EnvelopeRecordUpdate) SetVersion(v string) { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s EnvelopeRecordUpdate) VersionPrev() string { return C.Struct(s).GetObject(2).ToText() }
func (s EnvelopeRecordUpdate) VersionPrevBytes() []byte {
	return C.Struct(s).GetObject(2).ToDataTrimLastByte()
}
func (s EnvelopeRecordUpdate) SetVersionPrev(v string) { C.Struct(s).SetObject(2, s.Segment.NewText(v)) }
func (s EnvelopeRecordUpdate) WriteJSON(w io.Writer) error {
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
	_, err = b.WriteString("\"versionPrev\":")
	if err != nil {
		return err
	}
	{
		s := s.VersionPrev()
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
func (s EnvelopeRecordUpdate) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s EnvelopeRecordUpdate) WriteCapLit(w io.Writer) error {
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
	_, err = b.WriteString("versionPrev = ")
	if err != nil {
		return err
	}
	{
		s := s.VersionPrev()
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
func (s EnvelopeRecordUpdate) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type EnvelopeRecordUpdate_List C.PointerList

func NewEnvelopeRecordUpdateList(s *C.Segment, sz int) EnvelopeRecordUpdate_List {
	return EnvelopeRecordUpdate_List(s.NewCompositeList(0, 3, sz))
}
func (s EnvelopeRecordUpdate_List) Len() int { return C.PointerList(s).Len() }
func (s EnvelopeRecordUpdate_List) At(i int) EnvelopeRecordUpdate {
	return EnvelopeRecordUpdate(C.PointerList(s).At(i).ToStruct())
}
func (s EnvelopeRecordUpdate_List) ToArray() []EnvelopeRecordUpdate {
	n := s.Len()
	a := make([]EnvelopeRecordUpdate, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s EnvelopeRecordUpdate_List) Set(i int, item EnvelopeRecordUpdate) {
	C.PointerList(s).Set(i, C.Object(item))
}
