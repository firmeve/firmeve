package main

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel"
	"github.com/pkg/errors"
)

func main() {
	error2()
}

func error2()  {
	e := kernel.Error("sss")
	e2 := kernel.ErrorWarp(e)
	e3 := kernel.ErrorWarp(e2)
	fmt.Printf("%v",e3.StackString())
}

func errorsw() {
	e := errors.New("sss")
	e2 := errors.WithStack(e)
	e3 := errors.WithStack(e2)

	fmt.Printf("%+v",e3)
}