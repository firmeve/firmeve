package firmeve

import (
	"github.com/firmeve/firmeve/utils"
	"reflect"
	"strings"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
	mutex   sync.Mutex
)

type prototypeFunc func(container Container, params ...interface{}) interface{}

type Container interface {
	Get(abstract interface{}, params ...interface{}) interface{}
	Has(abstract interface{}) bool
}

type binding struct {
	share     bool
	instance  interface{}
	prototype prototypeFunc
}

type Firmeve struct {
	Container
	bindings map[reflect.Type]*binding
	aliases  map[string]reflect.Type
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[reflect.Type]*binding),
			aliases:  make(map[string]reflect.Type),
		}
	})

	return firmeve
}

type bindingOption struct {
	name      string
	share     bool
	cover     bool
	prototype interface{}
}

func (f *Firmeve) Get(abstract interface{}, params ...interface{}) interface{} {

	reflectType := reflect.TypeOf(abstract)

	kind := reflectType.Kind()

	// name 解析
	if kind == reflect.String {
		if abstractReflectType, ok := f.aliases[strings.ToLower(abstract.(string))]; ok {
			return f.parseBindingPrototype(f.bindings[abstractReflectType])
		} else {
			panic(`object that does not exist`)
		}
	} else if kind == reflect.Func {
		newParams := make([]reflect.Value, 0)
		for i := 0; i < reflectType.NumIn(); i++ {
			result := f.Get(reflect.ValueOf(reflectType.In(i)).Interface())
			newParams = append(newParams, reflect.ValueOf(result))
		}

		if reflectType.NumOut() == 1 {
			return reflect.ValueOf(abstract).Call(newParams)[0].Interface()
		} else {
			return reflect.ValueOf(abstract).Call(newParams)
		}
	} else if kind == reflect.Ptr {
		newReflectType := reflectType.Elem()
		if newReflectType.Kind() == reflect.Struct {
			return f.parseStruct(newReflectType, reflect.ValueOf(abstract).Elem()).Addr().Interface()
		}
	} else {
		if binding, ok := f.bindings[reflectType]; ok {
			return f.parseBindingPrototype(binding)
		}
	}

	panic(`unsupported type`)
}

func WithBindShare(share bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).share = share
	}
}

func WithBindName(name string) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).name = strings.ToLower(name)
	}
}

func WithBindInterface(object interface{}) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).prototype = object
	}
}

func WithBindCover(cover bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).cover = cover
	}
}

func (f *Firmeve) Bind(options ...utils.OptionFunc) { //, value interface{}
	// 参数解析 default bindingOption
	bindingOption := newBindingOption()
	utils.ApplyOption(bindingOption, options...)

	reflectType := reflect.TypeOf(bindingOption.prototype)

	// 别名获取
	if bindingOption.name == `` {
		bindingOption.name = f.parsePathName(reflectType)
	}

	// 覆盖检测
	if _, ok := f.bindings[reflectType]; ok && !bindingOption.cover {
		panic("binding alias type already exists")
	}
	if _, ok := f.aliases[bindingOption.name]; ok && !bindingOption.cover {
		panic("binding alias type already exists")
	}

	mutex.Lock()
	defer mutex.Unlock()

	f.aliases[bindingOption.name] = reflectType
	f.bindings[reflectType] = newBinding(bindingOption)
}

func (f *Firmeve) Has(abstract interface{}) bool {

	reflectType := reflect.TypeOf(abstract)

	if _, ok := f.bindings[reflectType]; ok {
		return true
	}

	if reflectType.Kind() == reflect.String {
		if _, ok := f.aliases[abstract.(string)]; ok {
			return true
		}
	}

	return false
}

func (f *Firmeve) All() map[reflect.Type]*binding {
	return f.bindings
}

func (f *Firmeve) Aliases() map[string]reflect.Type {
	return f.aliases
}

//func (f *Firmeve) ReflectType(id string) reflect.Type {
//	return f.aliases[id]
//}

// parse binding object prototype
func (f *Firmeve) parseBindingPrototype(binding *binding, params ...interface{}) interface{} {
	if binding.share && binding.instance != nil {
		return binding.instance
	}

	prototype := binding.prototype(f, params...)
	var result interface{}
	if reflect.TypeOf(prototype).Kind() == reflect.Func {
		result = f.Get(prototype)
	} else {
		result = prototype
	}

	if binding.share {
		binding.instance = result
	}

	return result
}

// parse struct fields and auto binding field
func (f *Firmeve) parseStruct(reflectType reflect.Type, reflectValue reflect.Value) reflect.Value {
	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("inject")
		if tag != `` && reflectValue.Field(i).CanSet() {
			reflectValue.Field(i).Set(reflect.ValueOf(f.Get(tag)))
		}
	}

	return reflectValue
}

// parse struct
func (f *Firmeve) parsePathName(reflectType reflect.Type) string {
	var name string

	kind := reflectType.Kind()
	switch kind {
	case reflect.Ptr:
		name = strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, `.`)
	default:
		name = reflectType.Name()
	}

	if name == `` {
		panic(`the path name is empty`)
	}

	return strings.ToLower(name)
}

// ---------------------------- bindingOption ------------------------

func newBindingOption() *bindingOption {
	return &bindingOption{share: false, cover: false,}
}

// ---------------------------- binding ------------------------

func newBinding(option *bindingOption) *binding {
	return &binding{
		prototype: func(container Container, params ...interface{}) interface{} {
			return option.prototype
		},
		share: option.share,
	}
}
