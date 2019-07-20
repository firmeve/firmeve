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
	Resolve(abstract interface{}, params ...interface{}) interface{}
	Has(name string) bool
	Get(name string) interface{}
}

type binding struct {
	share       bool
	instance    interface{}
	prototype   prototypeFunc
	reflectType reflect.Type
}

type Firmeve struct {
	Container
	bindings map[string]*binding
	aliases  map[string]reflect.Type
	types    map[reflect.Type]string
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[string]*binding),
			aliases:  make(map[string]reflect.Type),
			types:    make(map[reflect.Type]string),
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

func (f *Firmeve) Get(name string) interface{} {
	if !f.Has(name) {
		panic(`object that does not exist`)
	}

	return f.parseBindingPrototype(f.bindings[strings.ToLower(name)])
}

// @todo 这里现在变成纯解析的方法
// @todo 如果不保存reflectType 那么就找不到binding的对像
// @todo 但reflectType 还需要测试，可能会出现重复，name必须是惟一的，但不同的Name可能会对应相同的reflectType
// @todo 特别是在相同类型不同名称，强制覆盖时，肯定会发生
// @todo 除非使用二维map，每次遍历，但可能也不会太准确,相同的map会有先后级，不一定先遍历到的就是的
func (f *Firmeve) Resolve(abstract interface{}, params ...interface{}) interface{} {

	reflectType := reflect.TypeOf(abstract)

	kind := reflectType.Kind()

	// name 解析
	//if kind == reflect.String {
	//	return f.parseBindingPrototype(f.bindings[abstract.(string)])
	//} else
	if kind == reflect.Func {
		newParams := make([]reflect.Value, 0)
		if len(params) != 0 {
			for param := range params {
				newParams = append(newParams, reflect.ValueOf(param))
			}
		} else {
			for i := 0; i < reflectType.NumIn(); i++ {
				name := f.types[reflectType.In(i)]
				newParams = append(newParams, reflect.ValueOf(f.Get(name)))
			}
		}

		if reflectType.NumOut() == 1 {
			return reflect.ValueOf(abstract).Call(newParams)[0].Interface()
		} else {
			return reflect.ValueOf(abstract).Call(newParams)
		}
	} else if kind == reflect.Ptr {
		newReflectType := reflectType.Elem()
		if name, ok := f.types[newReflectType]; ok {
			return f.Get(name)
		} else if newReflectType.Kind() == reflect.Struct {
			return f.parseStruct(newReflectType, reflect.ValueOf(abstract).Elem()).Addr().Interface()
		}
	}
	//fmt.Println("======================")
	//fmt.Printf("%#v", abstract)
	//fmt.Println("======================")
	//else {
	//	if binding, ok := f.bindings[reflectType]; ok {
	//		return f.parseBindingPrototype(binding)
	//	}
	//}

	panic(`unsupported type`)
}

func WithBindShare(share bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).share = share
	}
}

//func WithBindName(name string) utils.OptionFunc {
//	return func(option interface{}) {
//		option.(*bindingOption).name = strings.ToLower(name)
//	}
//}

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

func (f *Firmeve) Bind(name string, options ...utils.OptionFunc) { //, value interface{}
	// 参数解析 default bindingOption
	bindingOption := newBindingOption(name)
	utils.ApplyOption(bindingOption, options...)

	//reflectType := reflect.TypeOf(bindingOption.prototype)

	// 别名获取
	//if bindingOption.name == `` {
	//	bindingOption.name = f.parsePathName(reflectType)
	//}

	// 覆盖检测
	//if _, ok := f.bindings[reflectType]; ok && !bindingOption.cover {
	//	panic("binding alias type already exists")
	//}
	if _, ok := f.bindings[bindingOption.name]; ok && !bindingOption.cover {
		panic(`binding alias type already exists`)
	}

	mutex.Lock()
	defer mutex.Unlock()

	f.bindings[bindingOption.name] = newBinding(bindingOption)
	// 如果有多个相同的类型，会覆盖使用最后一个
	f.types[reflect.TypeOf(bindingOption.prototype)] = bindingOption.name
	//f.bindings[reflectType] = newBinding(bindingOption)
}

func (f *Firmeve) Has(name string) bool {
	if _, ok := f.bindings[strings.ToLower(name)]; ok {
		return true
	}

	return false
}

//func (f *Firmeve) All() map[reflect.Type]*binding {
//	return f.bindings
//}
//
//func (f *Firmeve) Aliases() map[string]reflect.Type {
//	return f.aliases
//}

//func (f *Firmeve) ReflectType(id string) reflect.Type {
//	return f.aliases[id]
//}

// parse binding object prototype
func (f *Firmeve) parseBindingPrototype(binding *binding, params ...interface{}) interface{} {
	if binding.share && binding.instance != nil {
		return binding.instance
	}

	prototype := binding.prototype(f, params...)

	if binding.share {
		binding.instance = prototype
	}

	return prototype
	//prototype := binding.prototype(f, params...)
	//var result interface{}
	//if reflect.TypeOf(prototype).Kind() == reflect.Func {
	//	result = f.Get(prototype)
	//} else {
	//	result = prototype
	//}
	//
	//if binding.share {
	//	binding.instance = result
	//}
	//
	//return result
}

// parse struct fields and auto binding field
func (f *Firmeve) parseStruct(reflectType reflect.Type, reflectValue reflect.Value) reflect.Value {
	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("inject")
		if tag != `` && reflectValue.Field(i).CanSet() {
			if _, ok := f.bindings[tag]; ok {
				reflectValue.Field(i).Set(reflect.ValueOf(f.Get(tag)))
			}
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

func newBindingOption(name string) *bindingOption {
	return &bindingOption{share: false, cover: false, name: strings.ToLower(name)}
}

// ---------------------------- binding ------------------------

func newBinding(option *bindingOption) *binding {
	return &binding{
		prototype: func(container Container, params ...interface{}) interface{} {
			if reflect.TypeOf(option.prototype).Kind() == reflect.Func {
				return container.Resolve(option.prototype)
			}
			return option.prototype
		},
		share: option.share,
	}
}
