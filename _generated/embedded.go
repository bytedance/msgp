package _generated

//go:generate msgp

type E1 struct {
	A string
	B string
}

type E2 struct {
	A string
	*F
	*G
}

type F struct {
	B string
}

type G struct {
}
