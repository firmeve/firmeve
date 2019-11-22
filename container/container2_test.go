package container

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/firmeve/firmeve/testdata/funcs"

	"github.com/firmeve/firmeve/testdata/structs"
)

func TestBaseContainer_Bind_Get(t *testing.T) {
	c := New()
	c.Bind("nesting", &structs.Nesting{
		NId: 10,
	})
	assert.Equal(t, 10, c.Get("nesting").(*structs.Nesting).NId)
	assert.Equal(t, c.Get("nesting"), c.Get("nesting"))
	c.Bind("nesting_struct", structs.Nesting{
		NId:      10,
		NBoolean: true,
	})
	assert.Equal(t, 10, c.Get("nesting_struct").(structs.Nesting).NId)
	assert.Equal(t, true, c.Get("nesting_struct").(structs.Nesting).NBoolean)

	c.Bind("dynamic", func() *structs.Nesting {
		return &structs.Nesting{
			NId: 15,
		}
	})
	reflect.ValueOf(c).CallSlice()
	fmt.Println(fmt.Sprintf("%p", c.Get("dynamic")))
	fmt.Println(fmt.Sprintf("%p", c.Get("dynamic")))
	fmt.Println(fmt.Sprintf("%p", c.Get("dynamic")))
	assert.NotEqual(t, fmt.Sprintf("%p", c.Get("dynamic")), fmt.Sprintf("%p", c.Get("dynamic")))
}

func TestBaseContainer_Make_Remove_Flush(t *testing.T) {
	c := New()
	c.Bind("dynamic", func() *structs.Nesting {
		return &structs.Nesting{
			NId: 15,
		}
	})
	assert.IsType(t, &structs.Nesting{}, c.Make("dynamic"))
	v := func(id int) int {
		return id
	}
	assert.Equal(t, 0, c.Make(v))
	assert.Equal(t, 10, c.Make(v, 10))

	c.Bind("nesting", &structs.Nesting{
		NId: 10,
	})

	assert.IsType(t, &structs.Nesting{}, c.Make("nesting"))
	assert.Equal(t, fmt.Sprintf("%p", c.Get("nesting")), fmt.Sprintf("%p", c.Get("nesting")))
	assert.Equal(t, 10, c.Make(new(structs.Nesting)).(*structs.Nesting).NId)

	c.Bind("nesting_struct", structs.Nesting{
		NId:      10,
		NBoolean: true,
	})
	assert.Equal(t, 10, c.Make("nesting_struct").(structs.Nesting).NId)

	assert.Equal(t, 0, c.Make(structs.Sample{}).(structs.Sample).Id)
	assert.Equal(t, 0, c.Make(&structs.Sample{}).(*structs.Sample).Id)

	c.Remove("nesting")
	assert.Panics(t, func() {
		c.Get("nesting")
	})

	_, ok := c.instances[reflect.TypeOf(new(structs.Nesting))]
	assert.Equal(t, false, ok)

	c.Flush()
	assert.Equal(t, 0, len(c.instances))
	assert.Equal(t, 0, len(c.bindings))
}

func TestBaseContainer_Make_Panic(t *testing.T) {
	c := New()

	assert.Panics(t, func() {
		c.Make(123)
	})

	assert.Panics(t, func() {
		c.Bind("abc", 123)
	})

	assert.Panics(t, func() {
		mfn := func(str string) []string {
			return []string{str}
		}
		value := c.Make(funcs.MixedFunc, mfn).(string)
		assert.Equal(t, "", value)
	})
}

func TestWithCover(t *testing.T) {
	c := New()
	c.Bind("nesting", &structs.Nesting{
		NId: 10,
	})

	c.Bind("nesting", &structs.Nesting{
		NId: 10,
	}, WithCover(true))

	assert.Panics(t, func() {
		c.Bind("nesting", &structs.Nesting{
			NId: 10,
		}, WithCover(false))
	})
}

func TestBaseContainer_Resolve_Func_Sample(t *testing.T) {
	c := New()
	fn := funcs.Sample
	v1 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn)).(*structs.Sample)
	v2 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn)).(*structs.Sample)
	assert.NotEqual(t, fmt.Sprintf("%p", v1), fmt.Sprintf("%p", v2))
	assert.Equal(t, 15, v1.Id)
	assert.Equal(t, 15, v2.Id)
}

func TestBaseContainer_Resolve_Func_NormalSample(t *testing.T) {
	c := New()
	fn := funcs.NormalSample
	values := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn)).([]interface{})
	assert.Equal(t, values[0].(int), 0)
	assert.Equal(t, values[1].(string), "")

	values1 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), 1, "testing").([]interface{})
	assert.Equal(t, values1[0].(int), 1)
	assert.Equal(t, values1[1].(string), "testing")

	values2 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), 1).([]interface{})
	assert.Equal(t, values2[0].(int), 1)
	assert.Equal(t, values2[1].(string), "")

	values3 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), "testing").([]interface{})
	assert.Equal(t, values3[0].(int), 0)
	assert.Equal(t, values3[1].(string), "testing")
}

func TestBaseContainer_Resolve_Func_StructFunc(t *testing.T) {
	c := New()
	fn := funcs.StructFunc
	value := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn))
	assert.Equal(t, 0, value.(int))

	c.Bind("sample", &structs.Sample{
		Id: 10,
	}, WithShare(true))
	value2 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn))
	assert.Equal(t, 10, value2.(int))

	value3 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), structs.Nesting{
		NId: 15,
	})
	assert.Equal(t, 25, value3.(int))
}

func TestBaseContainer_Resolve_Func_StructFuncFunc(t *testing.T) {
	c := New()
	fn := funcs.StructFuncFunc
	value := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), func(nesting *structs.Nesting) int {
		return nesting.NId
	})
	assert.Equal(t, 0, value.(int))

	c.Bind("sample", &structs.Sample{
		Id: 10,
	}, WithShare(true))
	value2 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), 10, func(nesting *structs.Nesting) int {
		return nesting.NId
	})
	assert.Equal(t, 20, value2.(int))

	value3 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), 10, &structs.Nesting{
		NId: 15,
	}, func(nesting *structs.Nesting) int {
		return nesting.NId
	})
	assert.Equal(t, 35, value3.(int))
}

func TestBaseContainer_Resolve_Func_MultipleParamSample(t *testing.T) {
	c := New()
	fn := funcs.MultipleParamSample
	values := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn)).([]interface{})
	assert.Equal(t, reflect.Ptr, reflect.TypeOf(values[0]).Kind())
	assert.Equal(t, reflect.Struct, reflect.TypeOf(values[0]).Elem().Kind())
	assert.Equal(t, "Sample", reflect.TypeOf(values[0]).Elem().Name())

	assert.Equal(t, reflect.Ptr, reflect.TypeOf(values[1]).Kind())
	assert.Equal(t, reflect.Struct, reflect.TypeOf(values[1]).Elem().Kind())
	assert.Equal(t, "Nesting", reflect.TypeOf(values[1]).Elem().Name())

	assert.Equal(t, reflect.Struct, reflect.TypeOf(values[2]).Kind())
	assert.Equal(t, "Nesting", reflect.TypeOf(values[2]).Name())

	assert.Equal(t, reflect.Struct, reflect.TypeOf(values[3]).Kind())
	assert.Equal(t, "Sample", reflect.TypeOf(values[3]).Name())
}

func TestBaseContainer_Resolve_Func_MixedFunc(t *testing.T) {
	c := New()
	fn := funcs.MixedFunc
	var mfn funcs.MapFunc
	mfn = func(str string) []string {
		return []string{str}
	}
	value := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), mfn).(string)
	assert.Equal(t, "", value)

	var slice2 funcs.TestSlice
	slice2 = []string{"abc"}
	value1 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), slice2, "abc", mfn).(string)
	assert.Equal(t, "abcabc", value1)
}

func TestBaseContainer_Resolve_Func_Func(t *testing.T) {
	c := New()
	fn := funcs.ParamFunc
	v1 := c.resolveFunc(reflect.TypeOf(fn), reflect.ValueOf(fn), 10, func(prt *structs.Sample) (*structs.Sample, int) {
		prt.Id = 20
		return prt, 10
	}).(*structs.Sample)
	fmt.Printf("%#v", v1.Id)
}

func TestBaseContainer_Resolve_Struct_Sample(t *testing.T) {
	c := New()
	value := c.resolveStruct2(reflect.TypeOf(new(structs.Sample))).(*structs.Sample)
	assert.Equal(t, reflect.Ptr, reflect.ValueOf(value).Kind())
	assert.Equal(t, 0, value.Id)
	assert.Equal(t, "", value.Title)
	assert.Equal(t, false, value.Boolean)
	assert.Equal(t, [3]int{0, 0, 0}, value.Array)
	assert.Equal(t, map[string]string{}, value.MapString)
	assert.Equal(t, []string{}, value.Slice)

	value1 := c.resolveStruct2(reflect.TypeOf(structs.Sample{})).(structs.Sample)
	assert.Equal(t, reflect.Struct, reflect.ValueOf(value1).Kind())
	assert.Equal(t, 0, value1.Id)
	assert.Equal(t, "", value1.Title)
	assert.Equal(t, false, value1.Boolean)
	assert.Equal(t, [3]int{0, 0, 0}, value1.Array)
	assert.Equal(t, map[string]string{}, value1.MapString)
	assert.Equal(t, []string{}, value1.Slice)
}

func TestBaseContainer_Resolve_Struct_Sample_NestingPrt_Bind(t *testing.T) {
	c := New()
	nesting := &structs.Nesting{
		NId: 15,
	}
	c.Bind("nesting", nesting, WithShare(true))

	value := c.resolveStruct2(reflect.TypeOf(new(structs.SampleNestingInject))).(*structs.SampleNestingInject)
	assert.Equal(t, true, reflect.DeepEqual(value.Nesting, value.Nesting2))
	assert.Equal(t, fmt.Sprintf("%p", value.Nesting), fmt.Sprintf("%p", value.Nesting2))
}

func TestBaseContainer_Resolve_Struct_Sample_NestingPrt(t *testing.T) {
	c := New()
	value := c.resolveStruct2(reflect.TypeOf(new(structs.SampleNestingPtr))).(*structs.SampleNestingPtr)
	assert.Equal(t, 0, value.Id)
	assert.Equal(t, "", value.Title)
	assert.Equal(t, false, value.Boolean)
	assert.Equal(t, [3]int{0, 0, 0}, value.Array)
	assert.Equal(t, map[string]string{}, value.MapString)
	assert.Equal(t, []string{}, value.Slice)
	assert.Equal(t, reflect.Ptr, reflect.ValueOf(value.Nesting).Kind())
	assert.Equal(t, 0, value.Nesting.NId)
	assert.Equal(t, "", value.Nesting.NTitle)
	assert.Equal(t, false, value.Nesting.NBoolean)
	assert.Equal(t, [3]int{0, 0, 0}, value.Nesting.NArray)
	assert.Equal(t, map[string]string{}, value.Nesting.NMapString)
	assert.Equal(t, []string{}, value.Nesting.NSlice)

	value1 := c.resolveStruct2(reflect.TypeOf(structs.SampleNesting{})).(structs.SampleNesting)
	assert.Equal(t, 0, value1.Id)
	assert.Equal(t, "", value1.Title)
	assert.Equal(t, false, value1.Boolean)
	assert.Equal(t, [3]int{0, 0, 0}, value1.Array)
	assert.Equal(t, map[string]string{}, value1.MapString)
	assert.Equal(t, []string{}, value1.Slice)
	assert.Equal(t, reflect.Ptr, reflect.ValueOf(value1.Nesting).Kind())
	assert.Equal(t, 0, value1.Nesting.NId)
	assert.Equal(t, "", value1.Nesting.NTitle)
	assert.Equal(t, false, value1.Nesting.NBoolean)
	assert.Equal(t, [3]int{0, 0, 0}, value1.Nesting.NArray)
	assert.Equal(t, map[string]string{}, value1.Nesting.NMapString)
	assert.Equal(t, []string{}, value1.Nesting.NSlice)
	assert.Equal(t, reflect.Struct, reflect.ValueOf(value1.NestingValue).Kind())
	assert.Equal(t, 0, value1.NestingValue.NId)
	assert.Equal(t, "", value1.NestingValue.NTitle)
	assert.Equal(t, false, value1.NestingValue.NBoolean)
	assert.Equal(t, [3]int{0, 0, 0}, value1.NestingValue.NArray)
	assert.Equal(t, map[string]string{}, value1.NestingValue.NMapString)
	assert.Equal(t, []string{}, value1.NestingValue.NSlice)
}

//
//func TestBaseContainer_Resolve_struct(t *testing.T) {
//	c := New()
//	c.Bind("sub", &structs.Sub{
//		SubPublicKey: "SubPublicKey2...###",
//	})
//
//	//accout := new(Account)
//	//fmt.Println(c.resolveStruct2(reflect.TypeOf(accout)))
//	//
//	main := new(structs.Main)
//	//sub := new(structs.Sub)
//	m := c.resolveStruct2(reflect.TypeOf(main)).(*structs.Main)
//	fmt.Printf("%#v\n", m)
//	//fmt.Println("==================")
//	//fmt.Println(m.PrtSub.SubPublicKey)
//	//fmt.Println("==================")
//	//c.resolveStruct2(reflect.TypeOf(sub), reflect.ValueOf(sub))
//}
//
//func TestBaseContainer_Bind(t *testing.T) {
//	// 不支持静态类型的绑定
//	// types应该更改为resolveTypes，表示已解析的类型
//
//	// Bind
//	// bind("string","string") -> share
//	// bind("bool",bool) -> share
//	// bind("any",reflect.Kind < 动态类型) -> share
//	// bind("struct",struct object) -> share
//	// bind("ptr", ptr object) -> share
//	// bind("slice", slice) -> share
//	// bind("array", array) -> share array 是值
//	// bind("func" ,function(xxx,xxx)) -> no share
//	// 只要是非函数类型即可绑定为单例
//
//	// Resolve
//	// 必须要支持新对象的创建
//	// 必须要支持递归解析
//	// resolve("string") -> 表示直接读取容器key Get()
//	// 01 resolve(new(struct)|ptr) // ptr类型的struct
//	// 02 resolve(func) // func 是有多种类型的参数
//	// 03 resolve(slice|array|struct|map) 创建一个新的slice|array|struct
//	// 如果 反射类型在container中存在并且是singleton那么则返回已存在的类型
//
//	v1 := [2]int{0, 1}
//	v2 := [3]int{0, 1}
//	fmt.Println(reflect.TypeOf(v1) == reflect.TypeOf(v2))
//	fmt.Println("!!!!!!!!!!!!!!!!!!!")
//
//	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(a))
//	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(b))
//	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(c))
//	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(d))
//	fmt.Println("====")
//	fmt.Println(reflect.TypeOf("abc") == reflect.TypeOf("b"))
//
//	fmt.Println("########")
//	z := []int{0, 1}
//	fmt.Println(reflect.TypeOf(z).Elem().Kind())
//	//type Account struct {
//	//	Id     uint32
//	//	Name   string
//	//	Nested struct {
//	//		Age uint8
//	//	}
//	//}
//	//account := &Account{
//	//	Id: 10, Name: "jim",
//	//	Nested: struct{ Age uint8 }{Age: 20},
//	//}
//	str := "abc"
//	fmt.Println(reflect.TypeOf(str).Kind(), reflect.ValueOf(str).Kind())
//	fmt.Println("########")
//	m := struct {
//		id        int
//		substruct *struct {
//			title string
//		}
//	}{
//		id: 1,
//		substruct: &struct {
//			title string
//		}{
//			title: "abc",
//		},
//	}
//	fmt.Println(reflect.TypeOf(m).Field(1))
//
//	//var b, c interface{}
//	//var d int
//	//if b == nil {
//	//	fmt.Println("111")
//	//}
//	//if d == nil {
//	//	fmt.Println("111")
//	//}
//	//fmt.Println(b == c)
//}
//
//func a(a string) string {
//	return a
//}
//func c(c string) string {
//	return c
//}
//func d(c string, d int) string {
//	return c
//}
//func b(b int) int {
//	return b
//}

//func Test(p1,p2,p3)  {
//}
//// 普通调用
//p1 := ...
//p2 := ...
//p3 := ...
//Test(p1,p2,p3)
//
////ioc
//x.Resolve(Test)
//
//x.Resolve(Test) == Test(p1,p2,p3)
