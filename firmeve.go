package firmeve

import (
	"github.com/firmeve/firmeve/utils"
	"github.com/kataras/iris/core/errors"
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
	Has(id string) bool
	Resolve(abstract interface{}) interface{}
	//Bind(id interface{})
}

type binding struct {
	share     bool
	instance  interface{}
	prototype prototypeFunc
}

type Firmeve struct {
	Container
	bindings map[reflect.Type]*binding
	//aliases  map[string]string
	//aliases map[string]map[reflect.Kind]reflect.Type
	aliases map[string]reflect.Type
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
			bindings: make(map[reflect.Type]*binding),
			//aliases:  make(map[string]map[reflect.Kind]reflect.Type),
			aliases: make(map[string]reflect.Type),
			//aliases:  make(map[string]string),
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

type bindingOption struct {
	name  string
	share bool
	cover bool
	//aliases   []string
	//types   map[reflect.Type]string
	prototype interface{}
}

func (f *Firmeve) Resolve(abstract interface{}) interface{} {

	reflectType := reflect.TypeOf(abstract)

	if reflectType.Kind() == reflect.Func {
		panic(`Func`)
	}

	if _, ok := f.bindings[reflectType]; !ok {
		panic(`Error`)
	}

	binding := f.bindings[reflectType]
	if binding.share && binding.instance != nil {
		return binding.instance
	}

	prototype := binding.prototype(f)
	if binding.share {
		binding.instance = prototype
	}

	return prototype
	//
	//fmt.Printf("%#v",binding)
	//fmt.Println("===========================")
	//
	//kind := reflectType.Kind()
	//
	//switch kind {
	//// get name
	//case reflect.String:
	//	abstractString := abstract.(string)
	//	if abstractReflectType, ok := f.aliases[abstractString]; ok {
	//
	//		abstractBinding := f.bindings[abstractReflectType]
	//		if abstractBinding.share && abstractBinding.instance != nil {
	//			return abstractBinding.instance
	//		}
	//
	//		prototype := abstractBinding.prototype(f)
	//		if abstractBinding.share {
	//			abstractBinding.instance = prototype
	//		}
	//		return prototype
	//	} else {
	//		panic(`object that does not exist`)
	//	}
	//case reflect.Ptr:
	//	fmt.Println("GGGGGGGGGGGGGGGG")
	//	//path, _ := f.parsePathName(reflectType)
	//	return f.bindings[reflectType].prototype(f)
	//
	//case reflect.Func:
	//	// 反射参数
	//	//params := reflectType.NumIn()
	//	newParams := make([]reflect.Value, 0)
	//	for i := 0; i < reflectType.NumIn(); i++ {
	//		//reflectSubType := reflectType.In(0)
	//		//name, err := f.parseName(reflectSubType, "")
	//		//if err != nil {
	//		//	panic(err)
	//		//}
	//
	//		result := f.Get(reflect.ValueOf(reflectType.In(0).Kind()).Interface())
	//		newParams = append(newParams, reflect.ValueOf(result))
	//		fmt.Println("====================")
	//		fmt.Printf("%#v\n", result)
	//		fmt.Println("====================")
	//		//result := f.bindings[name].prototype(f)
	//		//if reflectSubType.Kind() == reflect.Func {
	//		//	// 参数暂时为空
	//		//	result = reflect.ValueOf(result).Call([]reflect.Value{})
	//		//} else {
	//		//
	//		//}
	//		//
	//		//newParams = append(newParams, reflect.ValueOf(result))
	//	}
	//
	//	return reflect.ValueOf(abstract).Call(newParams)[0].Interface()
	//
	//default:
	//	panic(`暂不支持`)
	//}

	return nil
}

func (f *Firmeve) parseBindingPrototype(binding *binding) interface{} {

	if binding.share && binding.instance != nil {
		return binding.instance
	}

	prototype := binding.prototype(f)
	if binding.share {
		binding.instance = prototype
	}

	return prototype
}

func (f *Firmeve) parseStruct(reflectType reflect.Type, reflectValue reflect.Value) reflect.Value {
	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("inject")
		if tag != `` && reflectValue.Field(i).CanSet() {
			reflectValue.Field(i).Set(reflect.ValueOf(f.Get(tag)))
		}
	}

	return reflectValue
}

func (f *Firmeve) Get(abstract interface{}, params ...interface{}) interface{} {

	reflectType := reflect.TypeOf(abstract)

	kind := reflectType.Kind()

	if binding, ok := f.bindings[reflectType]; ok {
		return f.parseBindingPrototype(binding)
	}

	// name 解析
	if kind == reflect.String {
		if abstractReflectType, ok := f.aliases[abstract.(string)]; ok {
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
	}
	//else if kind == reflect.Struct {
	//	reflectValue := reflect.ValueOf(abstract)
	//	for i := 0; i < reflectType.NumField(); i++ {
	//		tag := reflectType.Field(i).Tag.Get("inject")
	//		fmt.Println(tag,reflectValue.Field(i).CanSet())
	//		if tag != `` && reflectValue.Field(i).CanSet() {
	//			reflectValue.Field(i).Set(reflect.ValueOf(f.Get(tag)))
	//		}
	//	}
	//	//return f.parseStruct(reflectType, reflect.ValueOf(abstract)).Interface()
	//}

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

	// 默认字符串可调用名称
	if bindingOption.name == `` {
		pathName, err := f.parsePathName(reflectType)
		if err != nil {
			panic(err)
		}

		bindingOption.name = pathName
	}

	// Name覆盖检测
	if _, ok := f.aliases[bindingOption.name]; ok && !bindingOption.cover {
		panic("binding alias type already exists")
	}

	mutex.Lock()
	defer mutex.Unlock()

	f.aliases[bindingOption.name] = reflectType
	f.bindings[reflectType] = newBinding(f.setPrototype(bindingOption.prototype), bindingOption.share)
}

func (f *Firmeve) parsePathName(reflectType reflect.Type) (string, error) {
	var name string

	kind := reflectType.Kind()
	switch kind {
	case reflect.Ptr:
		name = strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, `.`)
	case reflect.Struct:
		name = strings.Join([]string{reflectType.PkgPath(), reflectType.Name()}, `.`)
	default:
		return ``, errors.New(`不支持的类型`)
	}

	return strings.ToLower(name), nil
}

func (f *Firmeve) setPrototype(prototype interface{}) prototypeFunc {
	return func(container Container, params ...interface{}) interface{} {
		return prototype
	}
}

func (f *Firmeve) Has(id string) bool {
	//panic("implement me")
	return false

}

// ---------------------------- bindingOption ------------------------

func newBindingOption() *bindingOption {
	return &bindingOption{share: false, cover: false,}
}

// ---------------------------- binding ------------------------

func newBinding(prototype prototypeFunc, share bool) *binding {
	return &binding{
		prototype: prototype,
		share:     share,
	}
}
