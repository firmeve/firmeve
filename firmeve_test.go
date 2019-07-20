package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/testdata"
	"log"
	"reflect"
	"time"

	//"github.com/firmeve/firmeve/testdata"
	"github.com/stretchr/testify/assert"
	//"log"
	//"math/rand"
	"testing"
)

type message string

func TestFirmeve_Resolve_String_Alias_Type(t *testing.T) {
	f := NewFirmeve()

	f.Bind("message", WithBindInterface(message("message bar")))
	f.Bind("string", WithBindInterface("string value"))

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
	f.Bind(t.Name()+"num", WithBindInterface(num), )
	f.Bind(t.Name()+"fnum", WithBindInterface(fnum), )
	f.Bind(t.Name()+"IntAliasNum", WithBindInterface(IntAliasNum), )

	assert.Equal(t, num, f.Get(t.Name() + "num").(int32))
	assert.Equal(t, fnum, f.Get(t.Name() + "fnum").(float32))
	assert.Equal(t, IntAliasNum, f.Get(t.Name() + "IntAliasNum").(IntAlias))
}

func TestFirmeve_Resolve_Bool(t *testing.T) {
	f := NewFirmeve()
	f.Bind("bool", WithBindInterface(false))

	assert.Equal(t, false, f.Get("bool").(bool))
}

func TestFirmeve_Resolve_Struct_Prt(t *testing.T) {
	t1 := testdata.NewT1()

	f := NewFirmeve()
	f.Bind(t.Name()+"t1", WithBindInterface(t1))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))
}

type FuncType func() int64

func TestFirmeve_Resolve_Func_Simple(t *testing.T) {
	f := NewFirmeve()

	var z FuncType
	z = func() int64 {
		return time.Now().UnixNano()
	}

	fmt.Println(z())
	fmt.Println(z())
	fmt.Println(z())
	//t1 := testdata.NewT1()
	f.Bind(t.Name()+"z", WithBindInterface(z))
	//assert.IsType(t, z, f.Get(t.Name()+"z"))

	result := f.Get(t.Name() + "z").(int64)
	log.Println(result)
}

//
func TestFirmeve_Resolve_Cover(t *testing.T) {
	t1 := testdata.NewT1()
	//t3 := testdata.NewT1Sturct()

	f := NewFirmeve()
	f.Bind(t.Name()+"t1", WithBindInterface(t1), )
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

	f.Bind(t.Name()+"t1", WithBindInterface(t1), WithBindCover(true))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

	assert.Panics(t, func() {
		f.Bind(t.Name()+"t1", WithBindInterface(t1), WithBindCover(false))
	}, "binding alias type already exists")
}

// 测试单例
func TestFirmeve_Singleton(t *testing.T) {
	//t1 := testdata.NewT1
	//t3 := testdata.NewT1Sturct()
	f := NewFirmeve()
	f.Bind(t.Name()+"t1.singleton", WithBindInterface(testdata.NewT1), WithBindShare(true))
	assert.Equal(t, fmt.Sprintf("%p", f.Get(t.Name()+"t1.singleton")), fmt.Sprintf("%p", f.Get(t.Name()+"t1.singleton")))

	f.Bind(t.Name()+"t2.prototype", WithBindInterface(testdata.NewT1), WithBindShare(false))
	assert.NotEqual(t, fmt.Sprintf("%p", f.Get(t.Name()+"t2.prototype")), fmt.Sprintf("%p", f.Get(t.Name()+"t2.prototype")))
}

func TestReflectType(t *testing.T)  {
	// 字符串
	s1 := "abc"
	s2 := "def"
	s3 := "abc"
	fmt.Println(reflect.TypeOf(s1) == reflect.TypeOf(s2))
	fmt.Println(reflect.ValueOf(s1) == reflect.ValueOf(s3))
	// number
	n1 := 20
	n2 := 40
	var n3 int64
	n3 =50
	fmt.Println(reflect.TypeOf(n1) == reflect.TypeOf(n2))
	fmt.Println(reflect.TypeOf(n1) == reflect.TypeOf(n3))
	// bool
	b1 := false
	b2 := true
	fmt.Println(reflect.TypeOf(b1) == reflect.TypeOf(b2))
	// struct
	t1 := testdata.NewT1Sturct()
	t2 := testdata.NewT1Sturct()
	t2.Name = "def"
	fmt.Println(reflect.TypeOf(t1) == reflect.TypeOf(t2))
	// *struct
	t3 := testdata.NewT1()
	t4 := testdata.NewT1()
	t5 := testdata.NewT22()
	fmt.Println(reflect.TypeOf(t3) == reflect.TypeOf(t4))
	fmt.Println(reflect.TypeOf(t3) == reflect.TypeOf(t5))
	// func
}

// 测试struct字段反射
func TestFirmeve_Resolve_Struct_Field(t *testing.T) {
	f := NewFirmeve()
	t1 := testdata.NewT1()
	f.Bind("t1", WithBindInterface(t1))

	fmt.Printf("%#v\n",f.Resolve(testdata.NewT2))

	t1struct := testdata.NewT1Sturct()
	f.Bind("t1struct", WithBindInterface(t1struct))
	fmt.Printf("%#v",f.Resolve(testdata.NewTStruct))
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
////
func TestFirmeve_Bind_Struct_Prt2(t *testing.T) {
	t1 := testdata.NewT1()
	//t2 := testdata.NewT2(t1)

	f := NewFirmeve()
	f.Bind(t.Name() + "t1",WithBindInterface(t1))

	t2 := new(testdata.T2)
	//fmt.Printf("%#v\n", t2)
	result := f.Resolve(t2)
	fmt.Printf("%#v\n", result)
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
