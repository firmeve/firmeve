package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/utils"
	"path/filepath"
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
	Has(name string) bool
	Get(name string) interface{}
	Bind(name string, prototype interface{}, options ...utils.OptionFunc)
	Resolve(abstract interface{}, params ...interface{}) interface{}
	Remove(name string)
}

type ServiceProvider interface {
	Register()
	Boot()
}

type binding struct {
	name        string
	share       bool
	instance    interface{}
	prototype   prototypeFunc
	reflectType reflect.Type
}

type Firmeve struct {
	bashPath         string
	bindings         map[string]*binding
	types            map[reflect.Type]string
	serviceProviders map[string]ServiceProvider
	booted           bool
}

type bindingOption struct {
	name      string
	share     bool
	cover     bool
	prototype interface{}
}

type firmeveOption struct {
	force bool
}

type FirmeveServiceProvider struct {
	Firmeve *Firmeve `inject:"firmeve"`
}

// Create a new firmeve container
func NewFirmeve(basePath string) *Firmeve {
	if firmeve != nil {
		return firmeve
	}
	basePath, err := filepath.Abs(basePath)
	if err != nil {
		panic(err.Error())
	}
	once.Do(func() {
		firmeve = &Firmeve{
			bashPath:         basePath,
			bindings:         make(map[string]*binding),
			types:            make(map[reflect.Type]string),
			serviceProviders: make(map[string]ServiceProvider),
			booted:           false,
		}

		firmeve.Bind(`firmeve`, firmeve, WithBindShare(true))
		firmeve.Bind(`config`, config.NewConfig(strings.Join([]string{basePath, `config`}, `/`)))
	})

	return firmeve
}

// Get a existing instance
func GetFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	panic(`firmeve not exists`)
}

// Start all service providers at once
func (f *Firmeve) Boot() {
	if f.booted {
		return
	}

	for _, provider := range f.serviceProviders {
		provider.Boot()
	}

	f.booted = true
}

// Register force param
func WithRegisterForce(force bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*firmeveOption).force = force
	}
}

// Register a service provider
func (f *Firmeve) Register(name string, provider ServiceProvider, options ...utils.OptionFunc) {
	// Parameter analysis
	firmeveOption := newFirmeveOption()
	utils.ApplyOption(firmeveOption, options...)

	if f.HasProvider(name) && !firmeveOption.force {
		return
	}

	f.registerProvider(name, provider)

	provider.Register()

	if f.booted {
		provider.Boot()
	}
}

// Add a service provider to serviceProviders map
func (f *Firmeve) registerProvider(name string, provider ServiceProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	f.serviceProviders[name] = provider
}

// Determine if the provider exists
func (f *Firmeve) HasProvider(name string) bool {
	if _, ok := f.serviceProviders[name]; ok {
		return ok
	}

	return false
}

// Get a if the provider exists
// If not found then panic
func (f *Firmeve) getProvider(name string) ServiceProvider {
	if !f.HasProvider(name) {
		panic(fmt.Sprintf("the %s service provider not exists", name))
	}

	return f.serviceProviders[name]
}

// Get application base path
func (f *Firmeve) GetBasePath() string {
	return f.bashPath
}

// Determine whether the specified name object is included in the container
func (f *Firmeve) Has(name string) bool {
	if _, ok := f.bindings[strings.ToLower(name)]; ok {
		return true
	}

	return false
}

// Get a object from container
func (f *Firmeve) Get(name string) interface{} {
	if !f.Has(name) {
		panic(`object that does not exist`)
	}

	return f.bindings[strings.ToLower(name)].resolvePrototype(f)
}

// Bind method `share` param
func WithBindShare(share bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).share = share
	}
}

// Bind method `cover` param
func WithBindCover(cover bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).cover = cover
	}
}

// Bind a object to container
func (f *Firmeve) Bind(name string, prototype interface{}, options ...utils.OptionFunc) { //, value interface{}
	// Parameter analysis
	bindingOption := newBindingOption(name, prototype)
	utils.ApplyOption(bindingOption, options...)

	// Coverage detection
	if _, ok := f.bindings[bindingOption.name]; ok && !bindingOption.cover {
		panic(`binding alias type already exists`)
	}

	// set binding item
	f.setBindingItem(newBinding(bindingOption))
}

// Parsing various objects
func (f *Firmeve) Resolve(abstract interface{}, params ...interface{}) interface{} {
	reflectType := reflect.TypeOf(abstract)
	kind := reflectType.Kind()

	if kind == reflect.Func {
		newParams := make([]reflect.Value, 0)
		if len(params) != 0 {
			for param := range params {
				newParams = append(newParams, reflect.ValueOf(param))
			}
		} else {
			for i := 0; i < reflectType.NumIn(); i++ {
				if name, ok := f.types[reflectType.In(i)]; ok {
					newParams = append(newParams, reflect.ValueOf(f.Get(name)))
				} else {
					panic(`unable to find reflection parameter`)
				}
			}
		}

		results := reflect.ValueOf(abstract).Call(newParams)

		resultsInterface := make([]interface{}, 0)
		for _, result := range results {
			resultInterface := result.Interface()
			if err, ok := resultInterface.(error); ok && err != nil {
				panic(err.Error())
			}

			resultsInterface = append(resultsInterface, resultInterface)
		}

		if reflectType.NumOut() == 1 {
			return resultsInterface[0]
		} else {
			return resultsInterface
		}
	} else if kind == reflect.Ptr {
		newReflectType := reflectType.Elem()
		if name, ok := f.types[newReflectType]; ok {
			return f.Get(name)
		} else if newReflectType.Kind() == reflect.Struct {
			return f.parseStruct(newReflectType, reflect.ValueOf(abstract).Elem()).Addr().Interface()
		}
	}

	panic(`unsupported type`)
}

func (f *Firmeve) Remove(name string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(f.bindings, name)

	for key, v := range f.types {
		if v == name {
			delete(f.types, key)
			break
		}
	}
}

// Set a item to types and bindings
func (f *Firmeve) setBindingItem(b *binding) {
	mutex.Lock()
	defer mutex.Unlock()

	// Set binding
	f.bindings[b.name] = b

	// Set type
	// Only support prt,struct and func type
	// No support string,float,int... scalar type
	originalKind := b.reflectType.Kind()
	if originalKind == reflect.Ptr || originalKind == reflect.Struct {
		f.types[b.reflectType] = b.name
	} else if originalKind == reflect.Func {
		// This is mainly used as a non-singleton type, using function execution, each time returning a different instance
		// When it is a function, parse the function, get the current real type, only support one parameter, the function must have only one return value
		f.types[reflect.TypeOf(b.resolvePrototype(f))] = b.name
	}
}

// Parse struct fields and auto binding field
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

// ---------------------------- firmeveOption ------------------------

func newFirmeveOption() *firmeveOption {
	return &firmeveOption{force: false}
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

	if b.reflectType.Kind() == reflect.Func {
		prototypeFunction = func(container Container, params ...interface{}) interface{} {
			return container.Resolve(prototype)
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
