package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/testdata"
	"testing"
)

func b() int {
	return 123
}

type testReject struct {
	Name string
	Age int

}

func TestFirmeve_Bind_Struct_Prt(t *testing.T) {
	t1 := testdata.NewT1()

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1),WithBindName("t1.prt"))

	result := f.Resolve(testdata.T2Call)
	//result := f.Resolve("t1.prt")
	fmt.Printf("%#v",result)
}

func TestFirmeve_Bind(t *testing.T) {
	//c:=b
	//println(c())
	//
	//

	//testReject1 := &testReject{"simon",30}

	firmeve := NewFirmeve()
	//t1 := NewT1()
	//firmeve.Bind(t1)

	//firmeve.Bind(WithBindShare(true),WithBindInterface(func() (string,int) {
	//	return `abc`,10
	//}),WithBindName("abc"))

	//firmeve.Bind(WithBindShare(true),WithBindInterface(testReject))
	//firmeve.Bind(WithBindShare(true),WithBindInterface(testReject{"simon",30}))
z:=[]string{"a","b"}
	firmeve.Bind(WithBindName("abcd"),WithBindInterface(z))

	//firmeve.Bind(func() interface{} {
	//	return NewT1()
	//},false)

	//fmt.Printf("%#v",firmeve.Resolve(demo.NewT2).(T2))
}