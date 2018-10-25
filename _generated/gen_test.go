package _generated

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/henrylee2cn/msgp/msgp"
)

// benchmark encoding a small, "fast" type.
// the point here is to see how much garbage
// is generated intrinsically by the encoding/
// decoding process as opposed to the nature
// of the struct.
func BenchmarkFastEncode(b *testing.B) {
	v := &TestFast{
		Lat:  40.12398,
		Long: -41.9082,
		Alt:  201.08290,
		Data: []byte("whaaaaargharbl"),
	}
	var buf bytes.Buffer
	msgp.Encode(&buf, v)
	en := msgp.NewWriter(msgp.Nowhere)
	b.SetBytes(int64(buf.Len()))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
}

// benchmark decoding a small, "fast" type.
// the point here is to see how much garbage
// is generated intrinsically by the encoding/
// decoding process as opposed to the nature
// of the struct.
func BenchmarkFastDecode(b *testing.B) {
	v := &TestFast{
		Lat:  40.12398,
		Long: -41.9082,
		Alt:  201.08290,
		Data: []byte("whaaaaargharbl"),
	}

	var buf bytes.Buffer
	msgp.Encode(&buf, v)
	dc := msgp.NewReader(msgp.NewEndlessReader(buf.Bytes(), b))
	b.SetBytes(int64(buf.Len()))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.DecodeMsg(dc)
	}
}

func (a *TestType) Equal(b *TestType) bool {
	if len(a.Els) != len(b.Els) {
		return false
	}
	if len(a.RuneSlice) != len(b.RuneSlice) {
		return false
	}
	if len(a.Slice1) != len(b.Slice1) {
		return false
	}
	if len(a.Slice2) != len(b.Slice2) {
		return false
	}
	e1a, e1b := a.Els, b.Els
	ra, rb := a.RuneSlice, b.RuneSlice
	s1a, s1b := a.Slice1, b.Slice1
	s2a, s2b := a.Slice2, b.Slice2
	if len(a.Els) == 0 {
		a.Els, b.Els = nil, nil
	}
	if len(a.RuneSlice) == 0 {
		a.RuneSlice, b.RuneSlice = nil, nil
	}
	if len(a.Slice1) == 0 {
		a.Slice1, b.Slice1 = nil, nil
	}
	if len(a.Slice2) == 0 {
		a.Slice2, b.Slice2 = nil, nil
	}
	// compare times, appended, then zero out those
	// fields, perform a DeepEqual, and restore them
	ta, tb := a.Time, b.Time
	if !ta.Equal(tb) {
		return false
	}
	aa, ab := a.Appended, b.Appended
	if !bytes.Equal(aa, ab) {
		return false
	}
	a.Appended, b.Appended = nil, nil
	a.Time, b.Time = time.Time{}, time.Time{}
	ok := reflect.DeepEqual(a, b)
	a.Time, b.Time = ta, tb
	a.Appended, b.Appended = aa, ab
	a.Els, b.Els = e1a, e1b
	a.RuneSlice, b.RuneSlice = ra, rb
	a.Slice1, b.Slice1 = s1a, s1b
	a.Slice2, b.Slice2 = s2a, s2b

	return ok
}

// This covers the following cases:
//  - Recursive types
//  - Non-builtin identifiers (and recursive types)
//  - time.Time
//  - map[string]string
//  - anonymous structs
//
func Test1EncodeDecode(t *testing.T) {
	f := 32.00
	s := []string{}
	tt := &TestType{
		F: &f,
		Els: map[string]string{
			"thing_one": "one",
			"thing_two": "two",
		},
		Obj: struct {
			ValueA string `msg:"value_a"`
			ValueB []byte `msg:"value_b"`
		}{
			ValueA: "here's the first inner value",
			ValueB: []byte("here's the second inner value"),
		},
		Child:    nil,
		Time:     time.Now(),
		Appended: msgp.Raw([]byte{}), // '[]'
		SlicePtr: &s,
	}

	var buf bytes.Buffer

	err := msgp.Encode(&buf, tt)
	if err != nil {
		t.Fatal(err)
	}

	tnew := new(TestType)

	err = msgp.Decode(&buf, tnew)
	if err != nil {
		t.Error(err)
	}

	if !tt.Equal(tnew) {
		t.Logf("in: %#v", tt)
		t.Logf("out: %#v", tnew)
		t.Fatal("objects not equal")
	}

	tanother := new(TestType)

	buf.Reset()
	msgp.Encode(&buf, tt)

	var left []byte
	left, err = tanother.UnmarshalMsg(buf.Bytes())
	if err != nil {
		t.Error(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left", len(left))
	}

	if !tt.Equal(tanother) {
		t.Logf("in: %v", tt)
		t.Logf("out: %v", tanother)
		t.Fatal("objects not equal")
	}
}

func TestIssue168(t *testing.T) {
	buf := bytes.Buffer{}
	test := TestObj{}

	msgp.Encode(&buf, &TestObj{ID1: "1", ID2: "2"})
	msgp.Decode(&buf, &test)

	if test.ID1 != "1" || test.ID2 != "2" {
		t.Fatalf("got back %+v", test)
	}
}
