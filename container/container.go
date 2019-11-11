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
	Resolve(abstract interface{}, params ...interface{}) interface{}
	Remove(name string)
	Flush()
}

type prototypeFunc func(container Container, params ...interface{}) interface{}
type bindingType map[string]*binding
type typesType map[reflect.Type]string

type baseContainer struct {
	bindings bindingType
	types    typesType
}

type binding struct {
	name        string
	share       bool
	instance    interface{}
	prototype   interface{}
	reflectType reflect.Type
}

type bindingOption struct {
	name      string
	share     bool
	cover     bool
	prototype interface{}
}

var (
	mutex sync.Mutex
)

// Create a new container instance
func New() *baseContainer {
	return &baseContainer{
		bindings: make(bindingType),
		types:    make(typesType),
	}
}

// Determine whether the specified name object is included in the container
func (c *baseContainer) Has(name string) bool {
	if _, ok := c.bindings[strings.ToLower(name)]; ok {
		return true
	}

	return false
}

// Get a object from container
func (c *baseContainer) Get(name string) interface{} {
	if !c.Has(name) {
		panic(fmt.Errorf("object[%s] that does not exist", name))
	}

	bind := c.bindings[strings.ToLower(name)]
	// 这里bind.instance != nil可能有问题,需要测试
	if bind.share && bind.instance != nil {
		return bind.instance
	}

	return c.Resolve(bind.prototype)
	//name := strings.ToLower(name)
	//if c.bindings {
	//
	//}
	//return c.bindings[strings.ToLower(name)].resolvePrototype(c)
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
	// Parameter analysis
	bindingOption := support.ApplyOption(newBindingOption(name, prototype), options...).(*bindingOption)

	// set binding item
	c.setBindingItem(newBinding(bindingOption), bindingOption.cover)
}

// Parsing various objects
func (c *baseContainer) Resolve(abstract interface{}, params ...interface{}) interface{} {
	reflectType := reflect.TypeOf(abstract)
	reflectValue := reflect.ValueOf(abstract)

	kind := reflect2.KindElemType(reflectType)

	if kind == reflect.Func {
		return c.resolveFunc(reflectType, reflectValue, params...)
	} else if kind == reflect.Ptr || kind == reflect.Struct {
		return c.resolveStruct(reflectType, reflect.Indirect(reflectValue))
	} else if kind == reflect.String {
		return c.Get(abstract.(string))
	}

	// bind exists instance
	//if name, ok := c.types[reflectType]; ok {
	//	bind := c.bindings[name]
	//	// 这里要测试
	//	if bind.share && bind.instance == nil {
	//		bind.instance = ...
	//	}
	//}

	panic(fmt.Errorf("unsupported type %T", abstract))
}

// Remove a binding
func (c *baseContainer) Remove(name string) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	name = strings.ToLower(name)

	delete(c.bindings, name)

	for key, v := range c.types {
		if v == name {
			delete(c.types, key)
			break
		}
	}
}

// Flush container
func (c *baseContainer) Flush() {
	c.bindings = make(bindingType, 0)
	c.types = make(typesType, 0)
}

// Set a item to types and bindings
func (c *baseContainer) setBindingItem(b *binding, cover bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// Coverage detection
	// Force cover
	// If is not force cover
	if !cover {
		var bindExists, typeExists bool
		_, bindExists = c.bindings[b.name]
		_, typeExists = c.types[b.reflectType]
		if bindExists || typeExists {
			panic(fmt.Errorf("binding alias type %s already exists", b.name))
		}
	}

	// Set binding
	c.bindings[b.name] = b
	c.types[b.reflectType] = b.name

	return
	// Set type
	// Only support prt,struct and func type
	// No support string,float,int... scalar type
	//originalKind := reflect2.KindElemType(b.reflectType)
	//if originalKind == reflect.Ptr || originalKind == reflect.Struct {
	//	c.types[b.reflectType] = b.name
	//} else if originalKind == reflect.Func {
	//	// This is mainly used as a non-singleton type, using function execution, each time returning a different instance
	//	// When it is a function, parse the function, get the current real type, only support one parameter, the function must have only one return value
	//	c.types[reflect.TypeOf(b.resolvePrototype(c))] = b.name
	//}
}

func (c *baseContainer) resolveFunc(reflectType reflect.Type, reflectValue reflect.Value, params ...interface{}) interface{} {
	if len(params) == 0 {
		params = reflect2.CallInParameterType(reflectType, func(i int, param reflect.Type) interface{} {
			if name, ok := c.types[param]; ok {
				return c.Get(name)
			} else {
				fmt.Println(reflect.Indirect(reflect.New(param)).Interface())
				return reflect.Indirect(reflect.New(param.Elem())).Interface()
				//return reflect2.InterfaceValue(param)
				fmt.Println(reflect.New(param).Elem().Interface())
				return reflect.New(param).Elem().Interface()
				//kindType := reflect2.KindElemType(param)
				//kindType := param.Kind()
				//if kindType == reflect.New() {
				//
				//}
				//return param.
				//panic(`unable to find reflection parameter`)
			}
		})
	}

	results := reflect2.CallFuncValue(reflectValue, params...)

	if reflectType.NumOut() == 1 {
		return results[0]
	}

	return results
}

// Resolve struct fields and auto binding field
func (c *baseContainer) resolveStruct2(reflectType reflect.Type, reflectValue reflect.Value) interface{} {
	reflectValue = reflect.Indirect(reflectValue)
	reflect2.CallFieldType(reflectType, func(i int, field reflect.StructField) interface{} {
		//tag := field.Tag.Get("inject")
		fieldValue := reflectValue.Field(i)
		valueKind := fieldValue.Kind()

		if reflect2.CanSetValue(fieldValue) {
			tag := field.Tag.Get("inject")
			if c.Has(tag) {
				fieldValue.Set(reflect.ValueOf(c.Get(tag)))
			} else if bindName, ok := c.types[field.Type]; ok {
				fieldValue.Set(reflect.ValueOf(c.Get(bindName)))
			} else {
				if valueKind == reflect.Slice {
					fieldValue.Set(reflect.MakeSlice(field.Type, 0, 0))
				} else if valueKind == reflect.Map {
					fieldValue.Set(reflect.MakeMap(field.Type))
				} else if valueKind == reflect.Struct || valueKind == reflect.Ptr {
					if valueKind == reflect.Ptr {
						fieldValue = reflect.New(field.Type.Elem())
					}
					c.resolveStruct2(field.Type, reflect.Indirect(fieldValue))
				}
			}

			//
		}
		//if tag != `` && reflect2.CanSetValue(fieldValue) {
		//	if _, ok := c.bindings[tag]; ok {
		//		result := c.Resolve(c.Get(tag))
		//		// Non-same type of direct skip
		//		if reflect2.KindElemType(reflect.TypeOf(result)) == reflect2.KindElemType(field.Type) {
		//			fieldValue.Set(reflect.ValueOf(result))
		//		}
		//	}
		//}

		return nil
	})

	return reflect2.InterfaceValue(reflectType, reflectValue)
}

//func (c *baseContainer) resolveStatic(reflectType reflect.Type, reflectValue reflect.Value) {
//
//}

// Resolve struct fields and auto binding field
func (c *baseContainer) resolveStruct(reflectType reflect.Type, reflectValue reflect.Value) interface{} {
	reflect2.CallFieldType(reflectType, func(i int, field reflect.StructField) interface{} {
		tag := field.Tag.Get("inject")
		fieldValue := reflectValue.Field(i)
		if tag != `` && reflect2.CanSetValue(fieldValue) {
			if _, ok := c.bindings[tag]; ok {
				result := c.Resolve(c.Get(tag))
				// Non-same type of direct skip
				if reflect2.KindElemType(reflect.TypeOf(result)) == reflect2.KindElemType(field.Type) {
					fieldValue.Set(reflect.ValueOf(result))
				}
			}
		}

		return nil
	})

	return reflect2.InterfaceValue(reflectType, reflectValue)
}

// ---------------------------- bindingOption ------------------------

// Create a new binding option struct
func newBindingOption(name string, prototype interface{}) *bindingOption {
	return &bindingOption{share: false, cover: false, name: strings.ToLower(name), prototype: prototype}
}

// ---------------------------- binding ------------------------

// Create a new binding struct
func newBinding(option *bindingOption) *binding {
	binding := &binding{
		name:        option.name,
		reflectType: reflect.TypeOf(option.prototype),
	}
	binding.share = binding.getShare(option.share)
	//binding.prototype = binding.getPrototypeFunc(option.prototype)
	binding.prototype = option.prototype

	return binding
}

// Get share, If type kind is not func type
func (b *binding) getShare(share bool) bool {
	// Func type disabled set share status
	if b.reflectType.Kind() == reflect.Func {
		b.share = false
	}

	return share
}

// Parse package prototypeFunc type
func (b *binding) getPrototypeFunc(prototype interface{}) prototypeFunc {
	var prototypeFunction prototypeFunc

	if reflect2.KindElemType(b.reflectType) == reflect.Func {
		prototypeFunction = func(container Container, params ...interface{}) interface{} {
			return container.Resolve(prototype, params...)
		}
	} else {
		prototypeFunction = func(container Container, params ...interface{}) interface{} {
			return prototype
		}
	}

	return prototypeFunction
}

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
