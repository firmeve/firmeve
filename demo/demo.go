package demo

import (
	"github.com/firmeve/firmeve/bak"
)

//var Bind *firmeve.Binding

type Demo struct {
	Title string
	At    int
}

func Test() interface{} {
	return &bak.Bak{
		Title: "abc",
		At:    10,
	}
}

func Test2() string {
	return "abc"
}
