package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/testdata"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
)

type message string

func TestFirmeve_Resolve_String_Alias_Type(t *testing.T) {
	f := NewFirmeve()

	f.Bind(WithBindInterface(message("message bar")), WithBindName("message"))
	f.Bind(WithBindInterface("string value"), WithBindName("string"))

	//assert.IsType(t,message(""),f.Get("message"))
	z := f.Get("message").(message)
	assert.Equal(t, "message bar", fmt.Sprintf("%s", z))
	assert.Equal(t, "string value", f.Get("string").(string))
}

type IntAlias int32

func TestFirmeve_Resolve_Number(t *testing.T) {
	var num int32
	num = 255
	var fnum float32
	fnum = 255.02
	var IntAliasNum IntAlias
	IntAliasNum = 80

	f := NewFirmeve()
	f.Bind(WithBindInterface(num), WithBindName(t.Name()+"num"))
	f.Bind(WithBindInterface(fnum), WithBindName(t.Name()+"fnum"))
	f.Bind(WithBindInterface(IntAliasNum), WithBindName(t.Name()+"IntAliasNum"))

	assert.Equal(t, num, f.Get(t.Name() + "num").(int32))
	assert.Equal(t, fnum, f.Get(t.Name() + "fnum").(float32))
	assert.Equal(t, IntAliasNum, f.Get(t.Name() + "IntAliasNum").(IntAlias))
}

func TestFirmeve_Resolve_Bool(t *testing.T) {
	f := NewFirmeve()
	f.Bind(WithBindInterface(false), WithBindName("bool"))

	assert.Equal(t, false, f.Get("bool").(bool))
}

func TestFirmeve_Resolve_Struct_Prt(t *testing.T) {
	t1 := testdata.NewT1()

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1), WithBindName(t.Name()+"t1"))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

}

type FuncType func() int

func TestFirmeve_Resolve_Func_Simple(t *testing.T) {
	f := NewFirmeve()

	var z FuncType
	z = func() int {
		return rand.Int()
	}

	//t1 := testdata.NewT1()
	f.Bind(WithBindInterface(z), WithBindName(t.Name()+"z"))
	assert.IsType(t, z, f.Get(t.Name()+"z"))

	result := f.Get(t.Name() + "z").(FuncType)
	log.Println(result())
}

func TestFirmeve_Resolve_Cover(t *testing.T) {
	t1 := testdata.NewT1()
	//t3 := testdata.NewT1Sturct()

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1), WithBindName(t.Name()+"t1"))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

	f.Bind(WithBindInterface(t1), WithBindName(t.Name()+"t1"), WithBindCover(true))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

	assert.Panics(t, func() {
		f.Bind(WithBindInterface(t1), WithBindName(t.Name()+"t1"), WithBindCover(false))
	}, "binding alias type already exists")
}

// 测试单例
func TestFirmeve_Singleton(t *testing.T) {
	//t1 := testdata.NewT1
	//t3 := testdata.NewT1Sturct()

	f := NewFirmeve()
	f.Bind(WithBindInterface(testdata.NewT1), WithBindName(t.Name()+"t1.singleton"), WithBindShare(true))
	fmt.Printf("%p\n",f.Get(t.Name()+"t1.singleton"))
	fmt.Printf("%p\n",f.Get(t.Name()+"t1.singleton"))
	f.Bind(WithBindInterface(testdata.NewT1), WithBindName(t.Name()+"t11.singleton"), WithBindShare(false))
	fmt.Printf("%p\n",f.Get(t.Name()+"t11.singleton"))
	fmt.Printf("%p\n",f.Get(t.Name()+"t11.singleton"))
}

//
//func TestFirmeve_Bind_Struct_Prt(t *testing.T) {
//	var i int8 = -1
//	fmt.Printf("%#v\n", *(*uint8)(unsafe.Pointer(&i)))
//	t1 := testdata.NewT1()
//	t2 := testdata.NewT1()
//	t3 := testdata.NewT1Sturct()
//	fmt.Printf("%p\n", t1)
//	fmt.Printf("%p\n", t2)
//	fmt.Printf("%p\n", reflect.TypeOf(t1))
//	fmt.Printf("%p\n", reflect.TypeOf(t2))
//	fmt.Printf("%p\n", reflect.TypeOf(t3))
//	fmt.Println("===================")
//	fmt.Println(reflect.TypeOf(t1) == reflect.TypeOf(t3))
//	fmt.Println("===================")
//	eface := (unsafe.Pointer(&t1))
//	fmt.Printf("%p", eface)
//
//	f := NewFirmeve()
//	f.Bind(WithBindInterface(t1), WithBindName("t1.prt"))
//
//	//result := f.Get(testdata.T2Call)
//	result := f.Get("t1.prt")
//	fmt.Printf("%#v", result)
//}
//
//
func TestFirmeve_Bind_Struct_Prt2(t *testing.T) {
	t1 := testdata.NewT1()
	//t2 := testdata.NewT2(t1)

	f := NewFirmeve()
	f.Bind(WithBindInterface(t1), WithBindName("t1"))

	//t2 := new(testdata.T2)
	//fmt.Printf("%#v\n", t2)
	//result := f.Get(t2)
	//fmt.Printf("%#v\n", result.(*testdata.T2))
	//result.(*testdata.T2).Age = 10
	//fmt.Printf("%#v\n", t2)

	//t4 := testdata.T2{}
	//result2 := f.Get(t4)
	//fmt.Printf("%#v\n", result2.(testdata.T2))
}

//
//func TestFirmeve_Bind(t *testing.T) {
//	//c:=b
//	//println(c())
//	//
//	//
//
//	//testReject1 := &testReject{"simon",30}
//
//	firmeve := NewFirmeve()
//	//t1 := NewT1()
//	//firmeve.Bind(t1)
//
//	//firmeve.Bind(WithBindShare(true),WithBindInterface(func() (string,int) {
//	//	return `abc`,10
//	//}),WithBindName("abc"))
//
//	//firmeve.Bind(WithBindShare(true),WithBindInterface(testReject))
//	//firmeve.Bind(WithBindShare(true),WithBindInterface(testReject{"simon",30}))
//	z := []string{"a", "b"}
//	firmeve.Bind(WithBindName("abcd"), WithBindInterface(z))
//
//	//firmeve.Bind(func() interface{} {
//	//	return NewT1()
//	//},false)
//
//	//fmt.Printf("%#v",firmeve.Resolve(demo.NewT2).(T2))
//}
