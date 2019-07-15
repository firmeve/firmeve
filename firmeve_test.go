package firmeve

import (
	"fmt"
	"testing"
)

type T1 struct {
	Name string
}

func NewT1() *T1 {
	return &T1{"Simon"}
}

type T2 struct {
	t1 *T1
	Age int
}

func NewT2(f *T1) T2  {
	return T2{t1:f,Age:10}
}


func TestFirmeve_Bind(t *testing.T) {
	firmeve := NewFirmeve()
	t1 := NewT1()
	firmeve.Bind(t1)

	fmt.Printf("%#v",firmeve.Resolve(NewT2).(T2))
}