package _generated

import (
	"bytes"
	"testing"

	"github.com/bytedance/msgp/msgp"
)

func TestEncodeE1DecodeE2(t *testing.T) {
	v := E1{
		H1: []string{"h1"},
		A:  "a",
		H2: []string{"h2"},
		J:  &GridView{"0000"},
		B:  "b",
	}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := E2{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}
	if !(vn.A == "a" &&
		vn.F.B == "b" && vn.G.F.B == "b" &&
		vn.G.H1[0] == "h1" && vn.G.H2[0] == "h2" &&
		vn.J.(*GridView).TestField == "0000") {
		t.Fail()
	}
	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func TestMarshalE1UnmarshalE2(t *testing.T) {
	v := E1{
		H1: []string{"h1"},
		A:  "a",
		H2: []string{"h2"},
		J:  &GridView{"0000"},
		B:  "b",
	}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	vn := E2{}
	left, err := vn.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if !(vn.A == "a" &&
		vn.F.B == "b" && vn.G.F.B == "b" &&
		vn.G.H1[0] == "h1" && vn.G.H2[0] == "h2" &&
		vn.J.(*GridView).TestField == "0000") {
		t.Fail()
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}
	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}
