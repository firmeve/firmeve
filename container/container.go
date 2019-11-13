package container

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/firmeve/firmeve/support"
	reflect2 "github.com/firmeve/firmeve/support/reflect"
)

type Container interface {
	Has(name string) bool
	Get(name string) interface{}
	Bind(name string, prototype interface{}, options ...support.Option)
	Make(abstract interface{}, params ...interface{}) interface{}
	Remove(name string)
	Flush()
}

type prototypeFunc func(container Container, params ...interface{}) interface{}
type bindingType map[string]*binding
type instanceType map[reflect.Type]interface{}

type baseContainer struct {
	bindings  bindingType
	instances instanceType
}

type binding struct {
	name        string
	prototype   prototypeFunc
	reflectType reflect.Type
}

type bindingOption struct {
	share bool
	cover bool
}

var (
	mutex sync.Mutex
)

// Create a new container instance
func New() *baseContainer {
	return &baseContainer{
		bindings:  make(bindingType),
		instances: make(instanceType),
	}
}

// Bind method `share` param
func WithShare(share bool) support.Option {
	return func(option support.Object) {
		option.(*bindingOption).share = share
	}
}

// Bind method `cover` param
func WithCover(cover bool) support.Option {
	return func(option support.Object) {
		option.(*bindingOption).cover = cover
	}
}

// Bind a object to container
func (c *baseContainer) Bind(name string, prototype interface{}, options ...support.Option) { //, value interface{}
	// 检测是否允许绑定
	reflectPrototypeType := reflect.TypeOf(prototype)
	// KindElemType => ptr elem 转换
	kind := reflect2.KindElemType(reflectPrototypeType)
	if kind <= reflect.Complex128 || kind == reflect.Chan || kind == reflect.UnsafePointer || kind == reflect.Func {
		panic(`only support array,slice,map,func,prt,struct`)
	}

	// Parameter analysis
	bindingOption := support.ApplyOption(&bindingOption{}, options...).(*bindingOption)

	_, ok := c.bindings[name]
	if !bindingOption.cover && ok {
		panic(fmt.Errorf("binding object %s already exists", name))
	}

	binding := &binding{
		name: name,
		prototype: func(container Container, params ...interface{}) interface{} {
			return container.Make(prototype, params...)
		},
		reflectType: reflectPrototypeType,
	}

	c.bindings[name] = binding

	// bind
	if bindingOption.share {
		c.instances[reflectPrototypeType] = c.resolvePrototype(binding)
	}
}

// Get a object from container
func (c *baseContainer) Get(name string) interface{} {
	if !c.Has(name) {
		panic(fmt.Errorf("object[%s] that does not exist", name))
	}

	bind := c.bindings[strings.ToLower(name)]
	// 如果是单例
	if v, ok := c.instances[bind.reflectType]; ok {
		return v
	}

	return c.resolvePrototype(bind)
}

//func (c *baseContainer) inTypeShare(reflectType reflect.Type) bool {
//	name, ok := c.types[reflectType]
//	if ok {
//		bind, bindOk := c.bindings[name]
//		if bindOk {
//			return bind.share
//		}
//	}
//	return false
//}

// Determine whether the specified name object is included in the container
func (c *baseContainer) Has(name string) bool {
	if _, ok := c.bindings[strings.ToLower(name)]; ok {
		return true
	}

	return false
}

// Determine whether the specified name object is included in the container
func (c *baseContainer) resolve(abstract interface{}, params ...interface{}) interface{} {
	reflectType := reflect.TypeOf(abstract)

	// 如果是单实例
	if _type, ok := c.instances[reflectType]; ok {
		return _type
	}

	kind := reflect2.KindElemType(reflectType)
	if kind == reflect.String && c.Has(abstract.(string)) {
		return c.Get(abstract.(string))
	} else if kind == reflect.Struct {
		return c.resolveStruct2(reflectType)
	} else if kind == reflect.Func {
		return c.resolveFunc(reflectType, reflect.ValueOf(abstract), params...)
	}

	panic(fmt.Errorf("unsupported type %T", abstract))
}

func (c *baseContainer) Make(abstract interface{}, params ...interface{}) interface{} {
	return c.resolve(abstract, params...)
}

// Remove a binding
func (c *baseContainer) Remove(name string) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	name = strings.ToLower(name)

	delete(c.instances, c.bindings[name].reflectType)
	delete(c.bindings, name)
}

// Flush container
func (c *baseContainer) Flush() {
	c.bindings = make(bindingType, 0)
	//c.instances = make(typesType, 0)
}

// Flush container
func (c *baseContainer) resolvePrototype(binding *binding) interface{} {
	return binding.prototype(c)
}

//
//// Set a item to types and bindings
//func (c *baseContainer) setBindingItem(b *binding, bindOption *bindingOption) {
//	mutex.Lock()
//	defer mutex.Unlock()
//
//	// Coverage detection
//	// Force cover
//	// If is not force cover
//	if !bindOption.cover {
//		var bindExists bool
//		//var bindExists, typeExists bool
//		_, bindExists = c.bindings[b.name]
//		if bindExists {
//			panic(fmt.Errorf("binding alias type %s already exists", b.name))
//		}
//		//_, typeExists = c.types[b.reflectType]
//		//if bindExists || typeExists {
//		//	panic(fmt.Errorf("binding alias type %s already exists", b.name))
//		//}
//	}
//
//	// Set binding
//	c.bindings[b.name] = b
//	//c.types[b.reflectType] = b.name
//
//	return
//	// Set type
//	// Only support prt,struct and func type
//	// No support string,float,int... scalar type
//	//originalKind := reflect2.KindElemType(b.reflectType)
//	//if originalKind == reflect.Ptr || originalKind == reflect.Struct {
//	//	c.types[b.reflectType] = b.name
//	//} else if originalKind == reflect.Func {
//	//	// This is mainly used as a non-singleton type, using function execution, each time returning a different instance
//	//	// When it is a function, parse the function, get the current real type, only support one parameter, the function must have only one return value
//	//	c.types[reflect.TypeOf(b.resolvePrototype(c))] = b.name
//	//}
//}

//func (c *baseContainer) resolveFunc2(reflectType reflect.Type, reflectValue reflect.Value, params ...interface{}) interface{} {
//	params = reflect2.CallInParameterType(reflectType, func(i int, param reflect.Type) interface{} {
//		//for _, givenParams := range params {
//		//
//		//}
//		if name, ok := c.types[param]; ok {
//			return c.Get(name)
//		} else {
//			fmt.Println(reflect.Indirect(reflect.New(param)).Interface())
//			return reflect.Indirect(reflect.New(param.Elem())).Interface()
//			//return reflect2.InterfaceValue(param)
//			fmt.Println(reflect.New(param).Elem().Interface())
//			return reflect.New(param).Elem().Interface()
//			//kindType := reflect2.KindElemType(param)
//			//kindType := param.Kind()
//			//if kindType == reflect.New() {
//			//
//			//}
//			//return param.
//			//panic(`unable to find reflection parameter`)
//		}
//	})
//
//	results := reflect2.CallFuncValue(reflectValue, params...)
//
//	if reflectType.NumOut() == 1 {
//		return results[0]
//	}
//
//	return results
//}
func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}

func (c *baseContainer) resolveFunc(reflectType reflect.Type, reflectValue reflect.Value, params ...interface{}) interface{} {
	// 合并参数
	reflectParams := reflect2.CallInParameterType(reflectType, func(i int, param reflect.Type) interface{} {
		// 查找已存在的参数匹配类型
		for _, inputParam := range params {
			inputParamType := reflect.TypeOf(inputParam)
			if param == inputParamType { // 如果参数类型相等
				return reflect.ValueOf(params[i])
			}
		}
		// 如果是单例参数
		if v, ok := c.instances[param]; ok {
			return reflect.ValueOf(v)
		}

		kind := param.Kind()
		if kind == reflect.Struct || kind == reflect.Ptr {
			return c.resolveStruct2(param)
		} else if kind == reflect.Func {
			// 这里有问题
			// 函数可能有多个返回值
			// reflect.New(reflect2.IndirectType(param)) 可能要 reflect.New(reflect2.IndirectType(param)).Elem()
			// @todo 暂时不支持多个返回值
			return c.resolveFunc(param, reflect.New(reflect2.IndirectType(param)))
		}

		return c.initZero(param)
	})

	results := reflect2.CallFuncValue(reflectValue, reflectParams...)

	if reflectType.NumOut() == 1 {
		return results[0]
	}

	return results
}

// Resolve struct fields and auto binding field
func (c *baseContainer) resolveStruct2(reflectType reflect.Type) interface{} {
	reflectValue := reflect.Indirect(reflect.New(reflect2.IndirectType(reflectType)))

	reflect2.CallFieldType(reflectType, func(i int, field reflect.StructField) interface{} {
		fieldValue := reflectValue.Field(i)
		kind := field.Type.Kind()
		if reflect2.CanSetValue(fieldValue) {
			tag := field.Tag.Get("inject")
			if c.Has(tag) {
				fieldValue.Set(reflect.ValueOf(c.Get(tag)))
			} else if v, ok := c.instances[field.Type]; ok {
				fieldValue.Set(reflect.ValueOf(v))
			} else if kind == reflect.Struct || kind == reflect.Ptr {
				fieldValue.Set(reflect.ValueOf(c.resolveStruct2(reflect2.IndirectType(field.Type))))
			} else {
				//fieldValue.Set(c.initZero(reflectType))
				fieldValue.Set(reflect.Zero(field.Type))
			}
		}

		return nil
	})

	return reflect2.InterfaceValue(reflectType, reflectValue)
}

func (c *baseContainer) initZero(reflectType reflect.Type) reflect.Value {
	kind := reflectType.Kind()
	//fmt.Println("===========")
	//fmt.Println(kind)
	//fmt.Println(reflectType.Elem().Name())
	//fmt.Println("===========")
	if kind == reflect.Slice {
		return reflect.MakeSlice(reflectType, 0, 0)
	} else if kind == reflect.Array {
		return reflect.New(reflect.ArrayOf(0, reflectType)).Elem()
	} else if kind == reflect.Map {
		return reflect.MakeMap(reflectType)
	} else if kind <= reflect.Complex128 {
		//return reflect.New(reflect2.IndirectType(reflectType)).Elem()
		return reflect.Zero(reflectType)
	}

	panic(fmt.Errorf("type error %s", kind))
}

//func (c *baseContainer) resolveStatic(reflectType reflect.Type, reflectValue reflect.Value) reflect.Value {
//	valueKind := reflectType.Kind()
//	if valueKind == reflect.Slice {
//		return reflect.MakeSlice(reflectType, 0, 0)
//	} else if valueKind == reflect.Map {
//		return reflect.MakeMap(reflectType)
//	} else if valueKind == reflect.Struct || valueKind == reflect.Ptr {
//		if valueKind == reflect.Ptr {
//			fieldValue = reflect.New(field.Type.Elem())
//		}
//		c.resolveStruct2(field.Type, reflect.Indirect(fieldValue))
//	}
//}

// Resolve struct fields and auto binding field
//func (c *baseContainer) resolveStruct(reflectType reflect.Type, reflectValue reflect.Value) interface{} {
//	reflect2.CallFieldType(reflectType, func(i int, field reflect.StructField) interface{} {
//		tag := field.Tag.Get("inject")
//		fieldValue := reflectValue.Field(i)
//		if tag != `` && reflect2.CanSetValue(fieldValue) {
//			if _, ok := c.bindings[tag]; ok {
//				result := c.Resolve(c.Get(tag))
//				// Non-same type of direct skip
//				if reflect2.KindElemType(reflect.TypeOf(result)) == reflect2.KindElemType(field.Type) {
//					fieldValue.Set(reflect.ValueOf(result))
//				}
//			}
//		}
//
//		return nil
//	})
//
//	return reflect2.InterfaceValue(reflectType, reflectValue)
//}

// ---------------------------- bindingOption ------------------------

//// Create a new binding option struct
//func newBindingOption(name string, prototype interface{}) *bindingOption {
//	return &bindingOption{share: false, cover: false, name: strings.ToLower(name), prototype: prototype}
//}

// ---------------------------- binding ------------------------
//
//// Create a new binding struct
//func newBinding(option *bindingOption) *binding {
//	binding := &binding{
//		name:        option.name,
//		reflectType: reflect.TypeOf(option.prototype),
//	}
//	binding.prototype = binding.prototypeFunc(option.prototype)
//	//binding.prototype = option.prototype
//
//	binding.share = binding.getShare(option.share)
//	//if binding.share {
//	//	binding.instance = container.Mak(option.prototype)
//	//}
//
//	return binding
//}
//
//// Get share, If type kind is not func type
//func (b *binding) getShare(share bool) bool {
//	// Func type disabled set share status
//	if b.reflectType.Kind() == reflect.Func {
//		b.share = false
//	}
//
//	return share
//}
//
//// Parse package prototypeFunc type
//func (b *binding) prototypeFunc(prototype interface{}) prototypeFunc {
//	return func(container Container, params ...interface{}) interface{} {
//		return container.Make(prototype, params...)
//	}
//	//
//	//var prototypeFunction prototypeFunc
//	//if b.reflectType.Kind() == reflect.Func {
//	//	prototypeFunction = func(container Container, params ...interface{}) interface{} {
//	//		return container.Make(prototype, params...)
//	//	}
//	//} else {
//	//	prototypeFunction = func(container Container, params ...interface{}) interface{} {
//	//		return prototype
//	//	}
//	//}
//	//
//	//return prototypeFunction
//}

//
//// Parse binding object prototype
//func (b *binding) resolvePrototype(container Container, params ...interface{}) interface{} {
//	if b.share && b.instance != nil {
//		return b.instance
//	}
//
//	prototype := b.prototype(container, params...)
//	if b.share {
//		b.instance = prototype
//	}
//
//	return prototype
//}
