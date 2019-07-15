package firmeve

import (
	"testing"
)


func TestFirmeve_Bind(t *testing.T) {
	firmeve := NewFirmeve()
	//t1 := NewT1()
	//firmeve.Bind(t1)

	firmeve.Bind(func() interface{} {
		return NewT1()
	},false)

	//fmt.Printf("%#v",firmeve.Resolve(demo.NewT2).(T2))
}