package reflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type (
	Foo struct {
		Id   int
		Name string
	}
)

func bar(n int, s string) {

}

func TestKindElemType(t *testing.T) {
	fPtr := new(Foo)
	v := KindElemType(reflect.TypeOf(fPtr))
	assert.Equal(t, v, reflect.Struct)

	foo := Foo{}
	v2 := KindElemType(reflect.TypeOf(foo))
	assert.Equal(t, v2, reflect.Struct)
}

func TestCallInParameterType(t *testing.T) {
	result := CallInParameterType(reflect.TypeOf(bar), func(i int, param reflect.Type) interface{} {
		return param.Kind()
	})
	assert.Equal(t, 2, len(result))
	assert.Equal(t, reflect.Int, result[0].(reflect.Kind))
	assert.Equal(t, reflect.String, result[1].(reflect.Kind))
}
