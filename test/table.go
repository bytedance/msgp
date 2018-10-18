package test

import "github.com/henrylee2cn/msgp/msgp"

//go:generate msgp

type Table struct {
	A msgp.Any
	B [5]msgp.Any
	C []msgp.Any
	D map[string]msgp.Any
	E map[string][]msgp.Any
}
