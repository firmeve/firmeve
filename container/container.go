package container

import (
	"fmt"
	"github.com/firmeve/firmeve/support"
	reflect2 "github.com/firmeve/firmeve/support/reflect"
	"reflect"
	"strings"
	"sync"
)

type Container interface {
	Has(name string) bool
	Get(name string) interface{}
	Bind(name string, prototype interface{}, options ...support.Option)
	Resolve(abstract interface{}, params ...interface{}) interface{}
	Remove(name string)
}

type BaseContainer struct {
	bindings map[string]*binding
	types    map[reflect.Type]string
}

type prototypeFunc func(container Container, params ...interface{}) interface{}

type binding struct {
	name        string
	share       bool
	instance    interface{}
	prototype   prototypeFunc
	reflectType reflect.Type
}

type bindingOption struct {
	name      string
	share     bool
	cover     bool
	prototype interface{}
}

var (
	container     *BaseContainer
	containerOnce sync.Once
)

// Create a new container instance
func NewContainer() *BaseContainer {
	if container != nil {
		return container
	}

	containerOnce.Do(func() {
		container = &BaseContainer{
			bindings: make(map[string]*binding),
			types:    make(map[reflect.Type]string),
		}
	})

	return container
}

// Determine whether the specified name object is included in the container
func (c *BaseContainer) Has(name string) bool {
	if _, ok := c.bindings[strings.ToLower(name)]; ok {
		return true
	}

	return false
}

// Get a object from container
func (c *BaseContainer) Get(name string) interface{} {
	if !c.Has(name) {
		panic(fmt.Errorf("object[%s] that does not exist", name))
	}

	return c.bindings[strings.ToLower(name)].resolvePrototype(c)
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
func (c *BaseContainer) Bind(name string, prototype interface{}, options ...support.Option) { //, value interface{}
	// Parameter analysis
	bindingOption := support.ApplyOption(newBindingOption(name, prototype), options...).(*bindingOption)

	// Coverage detection
	if _, ok := c.bindings[bindingOption.name]; ok && !bindingOption.cover {
		panic(fmt.Errorf("binding alias type %s already exists", bindingOption.name))
	}

	// set binding item
	c.setBindingItem(newBinding(bindingOption))
}

// Parsing various objects
func (c *BaseContainer) Resolve(abstract interface{}, params ...interface{}) interface{} {
	reflectType := reflect.TypeOf(abstract)
	reflectValue := reflect.ValueOf(abstract)

	kind := reflect2.KindType(reflectType)

	if kind == reflect.Func {
		return c.resolveFunc(reflectType, reflectValue, params...)
	} else if kind == reflect.Ptr || kind == reflect.Struct {
		return c.resolveStruct(reflect2.IndirectType(reflectType), reflect.Indirect(reflectValue))
	} else if kind == reflect.String {
		return c.Get(abstract.(string))
	}

	panic(fmt.Errorf("unsupported type %#v", abstract))
}

// Remove a binding
func (c *BaseContainer) Remove(name string) {
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

// Set a item to types and bindings
func (c *BaseContainer) setBindingItem(b *binding) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Set binding
	c.bindings[b.name] = b

	// Set type
	// Only support prt,struct and func type
	// No support string,float,int... scalar type
	originalKind := reflect2.KindType(b.reflectType)
	if originalKind == reflect.Ptr || originalKind == reflect.Struct {
		c.types[b.reflectType] = b.name
	} else if originalKind == reflect.Func {
		// This is mainly used as a non-singleton type, using function execution, each time returning a different instance
		// When it is a function, parse the function, get the current real type, only support one parameter, the function must have only one return value
		c.types[reflect.TypeOf(b.resolvePrototype(c))] = b.name
	}
}

func (c *BaseContainer) resolveFunc(reflectType reflect.Type, reflectValue reflect.Value, params ...interface{}) interface{} {
	if len(params) == 0 {
		params = reflect2.CallParameterType(reflectType, func(i int, param reflect.Type) interface{} {
			if name, ok := c.types[param]; ok {
				return c.Get(name)
			} else {
				panic(fmt.Errorf(`unable to find reflection parameter`))
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
func (c *BaseContainer) resolveStruct(reflectType reflect.Type, reflectValue reflect.Value) interface{} {
	reflect2.CallFieldType(reflectType, func(i int, field reflect.StructField) interface{} {
		tag := field.Tag.Get("inject")
		fieldValue := reflectValue.Field(i)
		if tag != `` && reflect2.CanSetValue(fieldValue) {
			if _, ok := c.bindings[tag]; ok {
				result := c.Resolve(c.Get(tag))
				// Non-same type of direct skip
				if reflect2.KindType(reflect.TypeOf(result)) == reflect2.KindType(field.Type) {
					fieldValue.Set(reflect.ValueOf(result))
				}
			}
		}

		return nil
	})

	return reflect2.InterfaceValue(reflectValue)
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
	binding.prototype = binding.getPrototypeFunc(option.prototype)

	return binding
}

// Get share, If type kind is not func type
func (b *binding) getShare(share bool) bool {
	if b.reflectType.Kind() != reflect.Func {
		b.share = true
	}

	return share
}

// Parse package prototypeFunc type
func (b *binding) getPrototypeFunc(prototype interface{}) prototypeFunc {
	var prototypeFunction prototypeFunc

	if reflect2.KindType(b.reflectType) == reflect.Func {
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

// Parse binding object prototype
func (b *binding) resolvePrototype(container Container, params ...interface{}) interface{} {
	if b.share && b.instance != nil {
		return b.instance
	}

	prototype := b.prototype(container, params...)
	if b.share {
		b.instance = prototype
	}

	return prototype
}