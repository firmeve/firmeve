package container

import (
	"fmt"
	"github.com/firmeve/firmeve/testdata"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
	"time"
)

type message string

func TestContainer_Resolve_String_Alias_Type(t *testing.T) {
	i := NewContainer()

	i.Bind("message", message("message bar"))
	i.Bind("string", "string value")

	//assert.IsType(t,message(""),f.Get("message"))
	z := i.Get("message").(message)
	assert.Equal(t, "message bar", fmt.Sprintf("%s", z))
	assert.Equal(t, "string value", i.Get("string").(string))
}

func TestBaseContainer_Has(t *testing.T) {
	i := NewContainer()
	i.Bind("string", "string value")
	assert.Equal(t, true, i.Has("string"))
}

type IntAlias int32
//
func TestContainer_Resolve_Number(t *testing.T) {
	var num int32
	num = 255
	var fnum float32
	fnum = 255.02
	var IntAliasNum IntAlias
	IntAliasNum = 80

	f := NewContainer()
	f.Bind(t.Name()+"num", num, )
	f.Bind(t.Name()+"fnum", fnum, )
	f.Bind(t.Name()+"IntAliasNum", IntAliasNum, )

	assert.Equal(t, num, f.Get(t.Name() + "num").(int32))
	assert.Equal(t, fnum, f.Get(t.Name() + "fnum").(float32))
	assert.Equal(t, IntAliasNum, f.Get(t.Name() + "IntAliasNum").(IntAlias))
}
//
func TestContainer_Resolve_Bool(t *testing.T) {
	f := NewContainer()
	f.Bind("bool", (false))

	assert.Equal(t, false, f.Get("bool").(bool))
}
//
func TestContainer_Resolve_Struct_Prt(t *testing.T) {
	t1 := testdata.NewT1()

	f := NewContainer()
	f.Bind(t.Name()+"t1", (t1))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))
}
//
type FuncType func() int64

func TestContainer_Resolve_Func_Simple(t *testing.T) {
	f := NewContainer()

	var z FuncType
	z = func() int64 {
		return time.Now().UnixNano()
	}

	fmt.Println(z())
	//t1 := testdata.NewT1()
	f.Bind(t.Name()+"z", z)
	//assert.IsType(t, z, f.Get(t.Name()+"z"))
	result := f.Get(t.Name() + "z").(int64)
	fmt.Println(result)
}
//
//
func TestContainer_Resolve_Cover(t *testing.T) {
	t1 := testdata.NewT1()
	//t3 := testdata.NewT1Sturct()

	f := NewContainer()
	f.Bind(t.Name()+"t1", (t1), )
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

	f.Bind(t.Name()+"t1", (t1), WithCover(true))
	assert.Equal(t, t1, f.Get(t.Name()+"t1"))

	assert.Panics(t, func() {
		f.Bind(t.Name()+"t1", (t1), WithCover(false))
	}, "binding alias type already exists")
}
//
// 测试单例
func TestContainer_Singleton(t *testing.T) {
	//t1 := testdata.NewT1
	//t3 := testdata.NewT1Sturct()
	f := NewContainer()
	f.Bind(t.Name()+"t1.singleton", (testdata.NewT1), WithShare(true))
	assert.Equal(t, fmt.Sprintf("%p", f.Get(t.Name()+"t1.singleton")), fmt.Sprintf("%p", f.Get(t.Name()+"t1.singleton")))

	f.Bind(t.Name()+"t2.prototype", (testdata.NewT1), WithShare(false))
	assert.NotEqual(t, fmt.Sprintf("%p", f.Get(t.Name()+"t2.prototype")), fmt.Sprintf("%p", f.Get(t.Name()+"t2.prototype")))
}

func TestReflectType(t *testing.T) {
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
	n3 = 50
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
//
// 测试struct字段反射
func TestContainer_Resolve_Struct_Field(t *testing.T) {
	f := NewContainer()
	t1 := testdata.NewT1()
	f.Bind(t.Name() + "t1", t1)

	fmt.Printf("%#v\n", f.Resolve(testdata.NewT2))

	t1struct := testdata.NewT1Sturct()
	f.Bind("t1struct", (t1struct))
	fmt.Printf("%#v", f.Resolve(testdata.NewTStruct))
}
//
// 测试非单例注入
func TestContainer_Resolve_Prototype(t *testing.T) {
	f := NewContainer()
	f.Bind("t1", (testdata.NewT1), WithCover(true))
	log.Printf("%#v\n", f.Resolve(testdata.NewT2))
	assert.IsType(t, testdata.NewT2(f.Get("t1").(*testdata.T1)), f.Resolve(testdata.NewT2))
}
//
//////
func TestContainer_Bind_Struct_Prt2(t *testing.T) {
	t1 := testdata.NewT1()

	f := NewContainer()
	f.Bind(t.Name()+"t1", (t1))

	t2 := new(testdata.T2)
	//fmt.Printf("%#v\n", t2)
	result := f.Resolve(t2)
	fmt.Printf("%#v\n", result)
}
//
func TestContainer_Remove(t *testing.T) {
	t1 := testdata.NewT1()
	//t2 := testdata.NewT2(t1)

	f := NewContainer()
	f.Bind(t.Name()+"t1", t1)
	f.Remove(t.Name() + "t1")

	assert.Equal(t, false, f.Has(t.Name()+"t1"))
}