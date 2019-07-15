package firmeve

import (
	"fmt"
	"go/importer"
	"reflect"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
)

type Container interface {
	Get(id string) interface{}
	Has(id string) bool
	//Bind(id interface{})
}

type binding struct {
	shared   bool
	instance interface{}
	object   func() interface{}
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

// 0. 每一个使用的都必须要注册，但可以不是单例，在init函数中是个好的选择
// 1. 全部存储函数类型
// 2. 如果是singleton保存到instance实例
// 3. 先实现这样的bind和resolve，后面增加struct tag 注入

func (f *Firmeve) Bind(object interface{}, share bool) { //, value interface{}
	//f.bindings[id] = newBinding(false, value)
	reflectType := reflect.TypeOf(object)
	println(reflectType.Kind().String())
	//if reflectType.Kind() == reflect.Ptr {
	//	fullName := strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, ".")
	//	f.bindings[fullName] = newBinding(share, object)
	//}

	//var method func(f *Firmeve) interface{}
	//
	//if reflectType.Kind() != reflect.Func {
	//	method = func(f *Firmeve) interface{} {
	//		return object
	//	}
	//} else {
	//	method = func(f *Firmeve) interface{} {
	//		reflect.ValueOf(object).Elem().Call()
	//	}
	//}



	if share {

	} else {

	}
}

func (f *Firmeve) Resolve(id interface{}) interface{} {
	reflectType := reflect.TypeOf(id)

	if reflectType.Kind() == reflect.Func {

		//params1Type := reflectType.In(0)//.Elem().Name()
		//z := nil
		////params1Value := &params1Type//reflect.NewAt(params1Type,unsafe.Pointer(params1Type.))
		//fmt.Printf("%#v",reflect.ValueOf(z).Convert(params1Type))

		// 重新设置一个新的struct对象，不过是nil
		//newObjPtr := reflect.New(params1Type)
		//newObj := reflect.Indirect(newObjPtr)

		//importer.Default()
		//fmt.Println(params1Type.Elem().PkgPath())
		pkg, err := importer.Default().Import(`github.com/firmeve/firmeve/demo`)
		if err != nil {
			fmt.Println(err.Error())
		}

		for _, declName := range pkg.Scope().Names() {
			fmt.Println(declName)
		}

		//println(params1Type,reflect.ValueOf(params1Value).Interface().(*T1))

		//fullName := strings.Join([]string{reflectType.In(0).Elem().PkgPath(), reflectType.In(0).Elem().Name()}, ".")
		//fmt.Println(fullName)
		//fmt.Printf("%#v",reflectType.In(0).Field(0))
		//stk := reflectType.In(0)
		//stkValue := reflect.ValueOf(reflect.TypeOf(id).In(0))
		////reflect.New(stk)
		////fmt.Printf("%#v",stkValue.MethodByName("NewT1").Call(make([]reflect.Value,0))[0].Interface())
		//fmt.Printf("%#v\n",stk.Elem().PkgPath())
		////zx := reflect.ValueOf(NewT1())
		//stkValue.Elem().Set(reflect.New(stk))
		//reflect.NewAt(stk, unsafe.Pointer(stk.UnsafeAddr())).Elem()
		//fmt.Printf("%#v",stkValue.Interface())
		//fmt.Printf("%#v",stkValue.Call([]reflect.Value{reflect.ValueOf("NewT1")}))

		//return reflect.ValueOf(id).Call([]reflect.Value{reflect.ValueOf(f.bindings[fullName].object)})[0].Interface()
	}

	return nil
}

// ---------------------------- binding ------------------------

func newBinding(shared bool, object func() interface{}) *binding {
	return &binding{
		shared:   shared,
		object:   object,
		instance: nil,
	}
}
