package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/utils"
	"go/importer"
	"reflect"
	"strings"
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
	shared    bool
	instance  interface{}
	prototype func(app Container, params ...interface{}) interface{}
}

type Firmeve struct {
	Container
	bindings map[string]*binding
	//bindingOptions
	//resolveOptions
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
// 0.0.0 如果是非singleton的，必须是func，类型，这样才能得到每次的原型
// 0.0 注册时候的路径包是个问题，没有路径包或名称就没法在resolve的时候进行取值
// 0.1 结构体可以通过反射取名称，但如果是func肯定是需要自己手动加入name
// 1. 全部存储函数类型
// 2. 如果是singleton保存到instance实例
// 3. 先实现这样的bind和resolve，后面增加struct tag 注入
// 4. 名称怎样实现惟一呢？，完整路径太长(github.com/b/c) 不是完整路径可能会存在冲突

type bindOption struct {
	name   string
	share  bool
	object interface{}
}

//type bindOptionFunc func(option *bindOption)

func WithBindShare(share bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindOption).share = share
	}
}

func WithBindName(name string) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindOption).name = strings.ToLower(name)
	}
}

func WithBindInterface(object interface{}) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindOption).object = object
	}
}

func (f *Firmeve) Bind(options ...utils.OptionFunc) { //, value interface{}

	//_,ok := Container.(f)
	//fmt.Println(ok)

	var bindOption = &bindOption{
		share: false,
	}
	// 参数解析
	utils.ApplyOption(bindOption, options...)

	// 处理对象
	object := bindOption.object
	reflectType := reflect.TypeOf(object)
	binding := newBinding()
	binding.shared = bindOption.share

	kind := reflectType.Kind()

	if kind == reflect.Func {
		if bindOption.name == "" {
			panic("函数类型，Name必须存在")
		}

		// 并且函数必须要有返回值
		//binding.
		//if reflectType.NumOut() == 0 {
		//	panic("函数类型，必须要有返回值")
		//}

		binding.prototype = func(app Container, params ...interface{}) interface{} {
			return object
		}

		f.bindings[bindOption.name] = binding

		func3 := f.bindings[bindOption.name].prototype(f)
		result := reflect.ValueOf(func3).Call([]reflect.Value{})
		fmt.Printf("%#v", result[1])
	} else if kind == reflect.Ptr {
		fmt.Println("============")
		// get Name

		name := strings.Join([]string{reflectType.Elem().PkgPath(),reflectType.Elem().Name()},`.`)
		fmt.Println(name)
		binding.prototype = func(container Container,params ...interface{}) interface{} {
			return object
		}

		f.bindings[name] = binding

		func3 := f.bindings[name].prototype(f)
		fmt.Printf("%#v",func3)
	} else if  kind == reflect.Struct  {
		name := strings.Join([]string{reflectType.PkgPath(),reflectType.Name()},`.`)
		fmt.Println(name,"==struct==")
		binding.prototype = func(container Container,params ...interface{}) interface{} {
			return object
		}

		f.bindings[name] = binding
		func3 := f.bindings[name].prototype(f)
		fmt.Printf("%#v",func3)
	}



	//fmt.Println(f.bindings[bindOption.name].prototype(f.(*Container)))
	//fmt.Printf("%#v",f.bindings)
	//fmt.Println(reflect.TypeOf(object).Kind().String())

	//if bindDefaultOption.name == `` {
	//
	//}

	//for _, opt := range options {
	//	opt(bindDefaultOption)
	//}

	//fmt.Printf("%#v", bindDefaultOption)

	//// 取得name和对象
	//paramsLen := len(params)
	//if paramsLen == 1 { //name需要自己反射，必须是结构体,share:false
	//
	//} else if paramsLen == 2 {
	//	if isShare, ok := params[1].(bool); ok && isShare { // name需要自己反射，单例类型
	//
	//	} else {
	//
	//	}
	//}
	//
	//fmt.Println(path.Base("a/b/c"))
	////f.bindings[id] = newBinding(false, value)
	//reflectType := reflect.TypeOf(object)
	//fmt.Printf("%p", reflectType)
	//reflectType.Elem()
	//println(reflectType.Kind().String())
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

	//if share {
	//
	//} else {
	//
	//}
}

func (f *Firmeve) Get(id string) interface{} {
	//panic("implement me")
	return "abc"
}

func (f *Firmeve) Has(id string) bool {
	//panic("implement me")
	return false

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

func newBinding() *binding {
	return &binding{
	}
}
