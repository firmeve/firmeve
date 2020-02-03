package container

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/firmeve/firmeve/support"
	reflect2 "github.com/firmeve/firmeve/support/reflect"
)

type (
	Container interface {
		Has(name string) bool
		Get(name string) interface{}
		Bind(name string, prototype interface{}, options ...support.Option)
		Make(abstract interface{}, params ...interface{}) interface{}
		Remove(name string)
		Flush()
	}

	baseContainer struct {
		bindings  bindingType
		instances instanceType
	}

	binding struct {
		name        string
		prototype   prototypeFunc
		reflectType reflect.Type
	}

	bindingOption struct {
		share bool
		cover bool
	}

	prototypeFunc func(container Container, params ...interface{}) interface{}

	bindingType map[string]*binding

	instanceType map[reflect.Type]interface{}
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
	if kind <= reflect.Complex128 || kind == reflect.Chan || kind == reflect.UnsafePointer {
		panic(`only support array,slice,map,func,prt(struct),struct`)
	}

	// Parameter analysis
	bindingOption := support.ApplyOption(&bindingOption{
		share: true,
		cover: false,
	}, options...).(*bindingOption)
	// func类型强制 force
	if kind == reflect.Func {
		bindingOption.share = false
	}

	_, ok := c.bindings[name]
	if !bindingOption.cover && ok {
		panic(fmt.Errorf("binding object %s already exists", name))
	}

	binding := &binding{
		name: name,
		prototype: func(container Container, params ...interface{}) interface{} {
			if kind == reflect.Func {
				return container.Make(prototype, params...)
			}

			return prototype
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
		panic(fmt.Errorf("object %s that does not exist", name))
	}

	binding := c.bindings[strings.ToLower(name)]
	// 如果是单例
	if v, ok := c.instances[binding.reflectType]; ok {
		return v
	}

	return c.resolvePrototype(binding)
}

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
	reflectValue := reflect.ValueOf(abstract)

	// 如果是单实例
	if _type, ok := c.instances[reflectType]; ok {
		return _type
	}

	kind := reflect2.KindElemType(reflectType)
	if kind == reflect.String && c.Has(abstract.(string)) {
		return c.Get(abstract.(string))
	} else if kind == reflect.Struct {
		return c.resolveStruct2(reflectType, reflectValue)
	} else if kind == reflect.Func {
		return c.resolveFunc(reflectType, reflectValue, params...)
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
	c.instances = make(instanceType, 0)
}

// Flush container
func (c *baseContainer) resolvePrototype(binding *binding) interface{} {
	return binding.prototype(c)
}

func (c *baseContainer) resolveFunc(reflectType reflect.Type, reflectValue reflect.Value, params ...interface{}) interface{} {
	// 合并参数
	reflectParams := reflect2.CallInParameterType(reflectType, func(i int, param reflect.Type) interface{} {
		// 查找已存在的参数匹配类型
		for _, inputParam := range params {
			inputParamType := reflect.TypeOf(inputParam)
			if param == inputParamType { // 如果参数类型相等
				return reflect.ValueOf(inputParam)
			}
		}
		// 如果是单例参数
		if v, ok := c.instances[param]; ok {
			return reflect.ValueOf(v)
		}

		kind := reflect2.KindElemType(param)
		// struct or ptr->struct
		if kind == reflect.Struct {
			return c.resolveStruct2(param, c.initZero(param))
		} else {
			return c.initZero(param)
		}
	})

	results := reflect2.CallFuncValue(reflectValue, reflectParams...)

	if reflectType.NumOut() == 1 {
		return results[0]
	}

	return results
}

// Resolve struct fields and auto binding field
func (c *baseContainer) resolveStruct2(reflectType reflect.Type, reflectValue2 reflect.Value) interface{} {
	if reflectValue2.IsZero() {
		reflectValue2 = c.initZero(reflectType)
	}

	reflectValue := reflect.Indirect(reflectValue2) //reflect.New(reflect2.IndirectType(reflectType)))

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
				fieldValue.Set(reflect.ValueOf(c.resolveStruct2(field.Type, fieldValue)))
			} else if fieldValue.IsZero() {
				fieldValue.Set(c.initZero(field.Type))
			} else {
				// nothing
			}
		}

		return nil
	})

	return reflect2.InterfaceValue(reflectType, reflectValue)
}

func (c *baseContainer) initZero(reflectType reflect.Type) reflect.Value {
	kind := reflect2.KindElemType(reflectType)
	if kind == reflect.Slice {
		return reflect.MakeSlice(reflectType, 0, 0)
	} else if kind == reflect.Array {
		return reflect.New(reflect.ArrayOf(reflectType.Len(), reflectType.Elem())).Elem()
	} else if kind == reflect.Map {
		return reflect.MakeMapWithSize(reflectType, 0)
	} else if kind == reflect.Struct {
		return reflect.New(reflect2.IndirectType(reflectType))
	} else if kind <= reflect.Complex128 || kind == reflect.String {
		return reflect.Zero(reflectType)
	}

	panic(fmt.Errorf("type error %s", kind))
}
