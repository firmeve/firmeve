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

func TestDefaultValueSample(t *testing.T) {
	sample := new(structs.SampleNesting)
	sample.Title = "testing title"
	sample.Id = 15
	//s := reflect.ValueOf(sample)
	//fmt.Println(s.IsZero())

	c := New()
	v := c.Make(sample).(*structs.SampleNesting)
	assert.Equal(t, 15, sample.Id)
	assert.Equal(t, "testing title", sample.Title)

	fmt.Printf("%#v\n", v)
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

func TestBaseContainer_Make(t *testing.T) {
	c := New()
	c.Bind(`number`, func(s string, type2 uint8) uint8 {
		return type2
	})
	v := c.resolve(`number`, "a", uint8(1))
	assert.Equal(t, uint8(1), v.(uint8))
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
	sample := new(structs.Sample)

	value := c.resolveStruct2(reflect.TypeOf(sample), reflect.ValueOf(sample)).(*structs.Sample)
	assert.Equal(t, reflect.Ptr, reflect.ValueOf(value).Kind())
	assert.Equal(t, 0, value.Id)
	assert.Equal(t, "", value.Title)
	assert.Equal(t, false, value.Boolean)
	assert.Equal(t, [3]int{0, 0, 0}, value.Array)
	assert.Equal(t, map[string]string{}, value.MapString)
	assert.Equal(t, []string{}, value.Slice)

	sample2 := structs.Sample{}
	value1 := c.resolveStruct2(reflect.TypeOf(sample2), reflect.ValueOf(sample2)).(structs.Sample)
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

	sn := new(structs.SampleNestingInject)
	value := c.resolveStruct2(reflect.TypeOf(sn), reflect.ValueOf(sn)).(*structs.SampleNestingInject)
	assert.Equal(t, true, reflect.DeepEqual(value.Nesting, value.Nesting2))
	assert.Equal(t, fmt.Sprintf("%p", value.Nesting), fmt.Sprintf("%p", value.Nesting2))
}

func TestBaseContainer_Resolve_Struct_Sample_NestingPrt(t *testing.T) {
	c := New()
	sn := new(structs.SampleNestingPtr)
	value := c.resolveStruct2(reflect.TypeOf(sn), reflect.ValueOf(sn)).(*structs.SampleNestingPtr)
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

	s3 := structs.SampleNesting{}
	value1 := c.resolveStruct2(reflect.TypeOf(s3), reflect.ValueOf(s3)).(structs.SampleNesting)
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
