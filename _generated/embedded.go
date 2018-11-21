package _generated

import "github.com/bytedance/msgp/msgp"

//go:generate msgp

type E1 struct {
	A  string
	B  string
	H1 []string
	H2 []string
	J  msgp.Any
}

type E2 struct {
	G
	A string
	*F
	J msgp.Any
}

type F struct {
	B string
}

type G struct {
	F
	H1 []string
	H2 []string
}
