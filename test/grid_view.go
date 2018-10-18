package test

import (
	"github.com/henrylee2cn/msgp/msgp"
)

//go:generate msgp

type GridView struct {
	TestField string
}

func init() {
	msgp.RegisterAny('g', new(GridView))
}
