package msgp

import (
	"fmt"
	"reflect"
	"sync"
)

// Any support any structure pointer to implement this interface.
type Any interface {
	Marshaler
	Unmarshaler
	Decodable
	Encodable
	Sizer
}

var (
	idToAnyType sync.Map // map[byte]reflect.Type
	anyTypeToId sync.Map // map[reflect.Type]byte
)

// RegisterAny records a type, identified by a value for that type, under its
// internal type id. That id will identify the concrete type of a value
// sent or received as an interface variable. Only types that will be
// transferred as implementations of interface values need to be registered.
// Expecting to be used only during initialization, it panics if the mapping
// between types and ids is not a bijection.
func RegisterAny(id byte, any Any) {
	if id == 0 {
		panic("attempt to register zero id")
	}
	if any == nil {
		panic("attempt to register nil")
	}

	t := reflect.Indirect(reflect.ValueOf(any)).Type()

	if nt, dup := idToAnyType.LoadOrStore(id, t); dup && nt != t {
		panic(fmt.Sprintf("msgp: registering duplicate types for %d: %v != %v", id, nt, t))
	}

	if nid, dup := anyTypeToId.LoadOrStore(t, id); dup && nid != id {
		idToAnyType.Delete(id)
		panic(fmt.Sprintf("msgp: registering duplicate ids for %v: %d != %d", t, nid, id))
	}
}

func EncodeAny(any Any, en *Writer) error {
	if any == nil {
		return en.WriteByte(0)
	}
	id, err := getAnyTypeId(any)
	if err != nil {
		return err
	}
	err = en.WriteByte(id)
	if err != nil {
		return err
	}
	return any.EncodeMsg(en)
}

func MarshalAny(any Any, o []byte) ([]byte, error) {
	if any == nil {
		return append(o, 0), nil
	}
	id, err := getAnyTypeId(any)
	if err != nil {
		return o, err
	}
	return any.MarshalMsg(append(o, id))
}

func DecodeAny(dc *Reader) (Any, error) {
	id, err := dc.ReadByte()
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	any, err := newAnyType(id)
	if err != nil {
		return nil, err
	}
	err = any.DecodeMsg(dc)
	if err != nil {
		return nil, err
	}
	return any, nil
}

func UnmarshalAny(bts []byte) (Any, []byte, error) {
	if len(bts) == 0 {
		return nil, bts, nil
	}
	id := bts[0]
	if id == 0 {
		return nil, bts[1:], nil
	}
	any, err := newAnyType(id)
	if err != nil {
		return nil, bts, err
	}
	bts, err = any.UnmarshalMsg(bts[1:])
	return any, bts, err
}

func Anysize(any Any) int {
	if any == nil {
		return 1
	}
	return any.Msgsize() + 1
}

func newAnyType(id byte) (Any, error) {
	t, ok := idToAnyType.Load(id)
	if !ok {
		return nil, fmt.Errorf("not support type ID: %d", id)
	}
	any, ok := reflect.New(t.(reflect.Type)).Interface().(Any)
	if !ok {
		return nil, fmt.Errorf("not support type ID: %d", id)
	}
	return any, nil
}

func getAnyTypeId(any Any) (byte, error) {
	t := reflect.Indirect(reflect.ValueOf(any)).Type()
	id, ok := anyTypeToId.Load(t)
	if !ok {
		return 0, fmt.Errorf("not support type ID: %d", id)
	}
	return id.(byte), nil
}
