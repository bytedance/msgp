package test

// Code generated by github.com/henrylee2cn/msgp DO NOT EDIT.

import (
	"github.com/henrylee2cn/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *GridView) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "TestField":
			z.TestField, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z GridView) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "TestField"
	err = en.Append(0x81, 0xa9, 0x54, 0x65, 0x73, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.TestField)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z GridView) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "TestField"
	o = append(o, 0x81, 0xa9, 0x54, 0x65, 0x73, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendString(o, z.TestField)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GridView) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "TestField":
			z.TestField, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z GridView) Msgsize() (s int) {
	s = 1 + 10 + msgp.StringPrefixSize + len(z.TestField)
	return
}
