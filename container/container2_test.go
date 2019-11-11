package container

import (
	"fmt"
	"github.com/firmeve/firmeve/testdata/structs"
	"reflect"
	"testing"
)

//type Account struct {
//	Id     uint32
//	Name   string
//	Nested *struct {
//		Age uint8
//	}
//}

func TestBaseContainer_Resolve_struct(t *testing.T) {
	c := New()
	c.Bind("sub", &structs.Sub{
		SubPublicKey: "SubPublicKey2",
	})

	main := new(structs.Main)
	//sub := new(structs.Sub)
	m := c.resolveStruct2(reflect.TypeOf(main), reflect.ValueOf(main))
	fmt.Printf("%#v", m)
	//c.resolveStruct2(reflect.TypeOf(sub), reflect.ValueOf(sub))
}

func TestBaseContainer_Bind(t *testing.T) {
	// Bind
	// bind("string","string") -> share
	// bind("bool",bool) -> share
	// bind("any",reflect.Kind < 动态类型) -> share
	// bind("struct",struct object) -> share
	// bind("ptr", ptr object) -> share
	// bind("slice", slice) -> share
	// bind("array", array) -> share array 是值
	// bind("func" ,function(xxx,xxx)) -> no share
	// 只要是非函数类型即可绑定为单例

	// Resolve
	// 必须要支持新对象的创建
	// 必须要支持递归解析
	// resolve("string") -> 表示直接读取容器key Get()
	// 01 resolve(new(struct)|ptr) // ptr类型的struct
	// 02 resolve(func) // func 是有多种类型的参数
	// 03 resolve(slice|array|struct) 创建一个新的slice|array|struct
	// 如果 反射类型在container中存在并且是singleton那么则返回已存在的类型

	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(a))
	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(b))
	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(c))
	fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(d))
	fmt.Println("====")
	fmt.Println(reflect.TypeOf("abc") == reflect.TypeOf("b"))

	fmt.Println("########")
	z := []int{0, 1}
	fmt.Println(reflect.TypeOf(z).Elem().Kind())
	//type Account struct {
	//	Id     uint32
	//	Name   string
	//	Nested struct {
	//		Age uint8
	//	}
	//}
	//account := &Account{
	//	Id: 10, Name: "jim",
	//	Nested: struct{ Age uint8 }{Age: 20},
	//}
	str := "abc"
	fmt.Println(reflect.TypeOf(str).Kind(), reflect.ValueOf(str).Kind())
	fmt.Println("########")
	m := struct {
		id        int
		substruct *struct {
			title string
		}
	}{
		id: 1,
		substruct: &struct {
			title string
		}{
			title: "abc",
		},
	}
	fmt.Println(reflect.TypeOf(m).Field(1))

	//var b, c interface{}
	//var d int
	//if b == nil {
	//	fmt.Println("111")
	//}
	//if d == nil {
	//	fmt.Println("111")
	//}
	//fmt.Println(b == c)
}

func a(a string) string {
	return a
}
func c(c string) string {
	return c
}
func d(c string, d int) string {
	return c
}
func b(b int) int {
	return b
}
