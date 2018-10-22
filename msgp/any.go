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

const anyHeaderLen = 2

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
func RegisterAny(namedID uint16, any Any) {
	if namedID == 0 {
		panic("attempt to register zero id")
	}
	if any == nil {
		panic("attempt to register nil")
	}

	t := reflect.Indirect(reflect.ValueOf(any)).Type()

	if nt, dup := idToAnyType.LoadOrStore(namedID, t); dup && nt != t {
		panic(fmt.Sprintf("msgp: registering duplicate types for %d: %v != %v", namedID, nt, t))
	}

	if nid, dup := anyTypeToId.LoadOrStore(t, namedID); dup && nid != namedID {
		idToAnyType.Delete(namedID)
		panic(fmt.Sprintf("msgp: registering duplicate ids for %v: %d != %d", t, nid, namedID))
	}
}

// EncodeAny encodes any struct pointer type.
func EncodeAny(any Any, en *Writer) error {
	cnt, err := createAnyBytes(any)
	if err != nil {
		return err
	}
	_, err = en.Write(AppendBytes(nil, cnt))
	return err
}

// MarshalAny marshals any struct pointer type.
func MarshalAny(any Any, o []byte) ([]byte, error) {
	cnt, err := createAnyBytes(any)
	if err != nil {
		return nil, err
	}
	o = AppendBytes(o, cnt)
	return o, nil
}

var scratchPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 512)
	},
}

// DecodeAny decodes any struct pointer type.
func DecodeAny(dc *Reader) (Any, error) {
	scratch := scratchPool.Get().([]byte)
	defer scratchPool.Put(scratch)
	scratch, err := dc.ReadBytes(scratch)
	if err != nil {
		return nil, err
	}
	any, err := parseAnyBytes(scratch)
	return any, err
}

// UnmarshalAny unmarshalls any struct pointer type.
func UnmarshalAny(bts []byte) (Any, []byte, error) {
	scratch := scratchPool.Get().([]byte)
	defer func() { scratchPool.Put(scratch) }()
	scratch, bts, err := ReadBytesBytes(bts, scratch)
	if err != nil {
		return nil, nil, err
	}
	any, err := parseAnyBytes(scratch)
	return any, bts, err
}

// Anysize returns the size of any struct pointer type.
func Anysize(any Any) int {
	if any == nil {
		return anyHeaderLen
	}
	return any.Msgsize() + anyHeaderLen
}

func newAnyType(id uint16) (Any, error) {
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

func getAnyTypeId(any Any) (uint16, error) {
	t := reflect.Indirect(reflect.ValueOf(any)).Type()
	id, ok := anyTypeToId.Load(t)
	if !ok {
		return 0, fmt.Errorf("not support type: %v", t)
	}
	return id.(uint16), nil
}

func createAnyBytes(any Any) ([]byte, error) {
	if any == nil {
		return nil, nil
	}
	id, err := getAnyTypeId(any)
	if err != nil {
		return nil, err
	}
	var tmp = make([]byte, anyHeaderLen)
	big.PutUint16(tmp, id)
	return any.MarshalMsg(tmp)
}

func parseAnyBytes(bts []byte) (Any, error) {
	if len(bts) < anyHeaderLen {
		return nil, nil
	}
	id := big.Uint16(bts)
	bts = bts[anyHeaderLen:]
	if id == 0 {
		return nil, nil
	}
	any, err := newAnyType(id)
	if err != nil {
		return nil, err
	}
	_, err = any.UnmarshalMsg(bts)
	return any, err
}
