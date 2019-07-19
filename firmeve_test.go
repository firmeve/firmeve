package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/testdata"
	"github.com/magiconair/properties/assert"
	"reflect"
	"testing"
	"unsafe"
)

func b() int {
	return 123
}

type testReject struct {
	Name string
	Age  int
}

type message string

func TestFirmeve_Resolve_String(t *testing.T) {
	f := NewFirmeve()

	f.Bind(WithBindInterface(message("abc")), WithBindName("message"))
	f.Bind(WithBindInterface("value"), WithBindName("value"))
	z := f.Get("message").(message)
	assert.Equal(t, "abc", fmt.Sprintf("%s",z))
	assert.Equal(t, "value", f.Get("value").(string))
}

func TestFirmeve_Resolve_Bool(t *testing.T) {
	f := NewFirmeve()
	f.Bind(WithBindInterface(false), WithBindName("bool"))

	assert.Equal(t, false, f.Get("bool").(bool))
}

func TestFirmeve_Resolve_Number(t *testing.T) {
	var num int32
	num = 255
	var fnum float32
	fnum = 255.02

	f := NewFirmeve()
	f.Bind(WithBindInterface(num), WithBindName("num"))
	f.Bind(WithBindInterface(fnum), WithBindName("fnum"))

	assert.Equal(t, num, f.Get("num").(int32))
	assert.Equal(t, fnum, f.Get("fnum").(float32))
}

func TestFirmeve_Resolve_Struct_Prt(t *testing.T) {
	t1 := testdata.NewT1()
	//t3 := testdata.NewT1Sturct()

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1), WithBindName("t1"))
	assert.Equal(t,t1,f.Get("t1"))

	f.Bind(WithBindInterface(t1), WithBindName("t3"))
	assert.Equal(t,t1,f.Get("t3"))
	//result := f.Get(testdata.T2Call)
	//result := f.Get("t1.prt")
	//fmt.Printf("%#v", result)
}

func TestFirmeve_Resolve_Cover(t *testing.T) {
	t1 := testdata.NewT1()
	//t3 := testdata.NewT1Sturct()

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1), WithBindName("t1"))
	assert.Equal(t,t1,f.Get("t1"))

	f.Bind(WithBindInterface(t1), WithBindName("t1"),WithBindCover(true))
	assert.Equal(t,t1,f.Get("t1"))

	assert.Panic(t, func() {
		f.Bind(WithBindInterface(t1), WithBindName("t1"),WithBindCover(false))
	},"binding alias type already exists")
}


func TestFirmeve_Bind_Struct_Prt(t *testing.T) {
	var i int8 = -1
	fmt.Printf("%#v\n", *(*uint8)(unsafe.Pointer(&i)))
	t1 := testdata.NewT1()
	t2 := testdata.NewT1()
	t3 := testdata.NewT1Sturct()
	fmt.Printf("%p\n", t1)
	fmt.Printf("%p\n", t2)
	fmt.Printf("%p\n", reflect.TypeOf(t1))
	fmt.Printf("%p\n", reflect.TypeOf(t2))
	fmt.Printf("%p\n", reflect.TypeOf(t3))
	fmt.Println("===================")
	fmt.Println(reflect.TypeOf(t1) == reflect.TypeOf(t3))
	fmt.Println("===================")
	eface := (unsafe.Pointer(&t1))
	fmt.Printf("%p", eface)

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1), WithBindName("t1.prt"))

	//result := f.Get(testdata.T2Call)
	result := f.Get("t1.prt")
	fmt.Printf("%#v", result)
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
	z := []string{"a", "b"}
	firmeve.Bind(WithBindName("abcd"), WithBindInterface(z))

	//firmeve.Bind(func() interface{} {
	//	return NewT1()
	//},false)

	//fmt.Printf("%#v",firmeve.Resolve(demo.NewT2).(T2))
}
