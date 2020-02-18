package a

import "github.com/firmeve/firmeve/contract/b"

type A interface {
	A_Demo(b b.B)
}

type S struct {
	b b.B
}

func (s *S) A_Demo(b b.B)  {
	s.b = b
}

func New() *S {
	return new(S)
}