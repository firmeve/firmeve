package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/tool"
	"reflect"
	"strings"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
	mutex   sync.Mutex
)

type Container interface {
	Get(id string) interface{}
	Has(id string) bool
}

type Binding struct {
	Shared bool
	Object interface{}
}

type binding struct {
	id       string
	shared   bool
	concrete func(object interface{}) func() interface{}
}

type Firmeve struct {
	bindings map[string]*binding
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[string]*binding),
		}
	})

	return firmeve
}

func (f *Firmeve) Bind(id string, concrete interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	f.bindings[id] = newBinding(id, false, concrete)
}

func (f *Firmeve) Has(id string) bool {
	_, ok := f.bindings[id]
	return ok
}

func (f *Firmeve) Resolve(id string, params ...interface{}) interface{} {
	//if f.Has(id) {
	//	object := f.bindings[id].object
	//} else {
	//	object := id
	//}
	//object := f.bindings[id].object
	//value := reflect.ValueOf(object)
	//z := make([]reflect.Value{},1)

	reflect.Copy(z,value)
	fmt.Printf("%#v",z)

	return nil
	//fmt.Println(reflect.TypeOf(object).PkgPath())
	//return reflect.ValueOf(object).Interface()
	//fmt.Printf("%s",reflect.ValueOf(object).Interface())
}

// ---------------------------- binding ------------------------

func newBinding(id string, shared bool, concrete interface{}) *binding {


	//fmt.Println(concreteType)

	return &binding{
		id:     id,
		shared: shared,
		concrete: func(concrete interface{}) func() interface{} {
			return func() interface{} {
				concreteType := reflect.TypeOf(concrete)
				elems := []interface{}{reflect.Ptr}
				if tool.InSlice(concreteType.Kind(), elems) {
					fmt.Println(concreteType.Elem().Name())
				} else {
					fmt.Println(concreteType.Name())
				}
				return concrete
			}
		},
	}
}


// 1. 主要的问题是结构体不能和php一样，直接写名字，这就有点尴尬
// 2. 强制写字符串，以包路径分隔可能也行，需要尝试，如github.com/b/c/d.x 形式
// 3. 如果结构体直接挂载对象，可以使用di来操作，di也需要new，递归反射,di别外引入，如下
// type A struct {
//	Db  *sql.DB `di:"db"`
//	Db1 *sql.DB `di:"db"`
//	B   *B      `di:"b,prototype"`
//	B1  *B      `di:"b,prototype"`
//}

// 4. 能否直接反射结构体实例化对象，必须是New(new)加结构体名称，如NewConfig或newConfig
// 5. 支持函数make,加上参数调用

func RelectTest(value interface{}, params ...interface{}) {

	//println()

	fmt.Println(reflect.ValueOf(value).String())

	if reflect.TypeOf(value).Name() == `string` {
		//name := reflect.ValueOf(value).String()
		name := reflect.ValueOf(value).String()
		slices := strings.Split(name, `;`)
		fmt.Println(slices[1])

		reflect.ValueOf(slices[1]).Call(make([]reflect.Value, 0))
	}

	//if reflect.TypeOf(value). {
	//	
	//}

	//fmt.Printf("%#v",reflect.TypeOf(value).PkgPath())
	//println(runtime.FuncForPC(reflect.ValueOf(value).Pointer()).Name())
	//
	//fmt.Printf("%s",reflect.TypeOf(value).Kind().String())
	//
	//result := reflect.ValueOf(value).Call(make([]reflect.Value,0))
	//fmt.Printf("%v",result)
	//fmt.Printf("%#v",result[0].Interface())

	//fmt.Printf("%v",reflect.TypeOf(value).())

	//fmt.Printf("%#v",value())

	//println(runtime.FuncForPC(reflect.ValueOf(value).Pointer()).Name())
	//fmt.Printf("%#v",reflect.TypeOf(value).Kind().String())
	//if reflect.TypeOf(value).Kind().String() == `func` {
	//	fmt.Println("GGGGGGGGGGGGG")
	//	v := reflect.ValueOf(value).Call(make([]reflect.Value,0))
	//	fmt.Println(v.(*bak.Bak))
	//
	//	//fmt.Printf("%#v",reflect.ValueOf(value).Call(make([]reflect.Value,0)))
	//}
	//getType := reflect.TypeOf(value)
	//value2 := reflect.ValueOf(value)
	//value
	//fmt.Printf("%#v",reflect.Typeof(value))

	//for i := 0; i < getType.NumField(); i++ {
	//	field := getType.Field(i)
	//	value := value2.Field(i).Interface()
	//	fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	//}

	//	fmt.Printf("%#v",value)
}
