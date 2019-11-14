package funcs

import (
	"github.com/firmeve/firmeve/testdata/structs"
)

func Sample() *structs.Sample {
	return &structs.Sample{
		Id: 15,
	}
}

func NormalSample(id int, str string) (int, string) {
	return id, str
}

func StructFunc(ptr *structs.Sample, nesting structs.Nesting) int {
	return ptr.Id + nesting.NId
}

func StructFuncFunc(id int, nesting *structs.Nesting, fn func(nesting *structs.Nesting) int, prt *structs.Sample) int {
	return id + fn(nesting) + prt.Id
}

type MapFunc func(str string) []string
type TestSlice []string

func MixedFunc(t TestSlice, str string, fn MapFunc) string {
	if len(t) > 0 {
		return t[0] + fn(str)[0]
	}

	return fn(str)[0]
}

func MultipleParamSample(ptr *structs.Sample, nestingPtr *structs.Nesting, nesting structs.Nesting, notPtr structs.Sample) (*structs.Sample, *structs.Nesting, structs.Nesting, structs.Sample) {
	return ptr, nestingPtr, nesting, notPtr
}

func ParamFunc(id int, ptr *structs.Sample, fn func(prt *structs.Sample) (*structs.Sample, int)) *structs.Sample {
	var ptr2 *structs.Sample
	var v int
	ptr2, v = fn(ptr)
	ptr2.Id = ptr2.Id + id + v
	return ptr2
}

//func c() {
//	c := &structs.Sample{}
//	ParamFunc(1, c, func(prt *structs.Sample) *structs.Sample {
//		return prt
//	})
//}
