package firmeve

import (
	"fmt"
	"sync"

	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/support"
)

const (
	Version               = "1.0.0"
	ModeDevelopment uint8 = iota
	ModeProduction
	ModeTesting
)

type Provider interface {
	Name() string
	Register()
	Boot()
}

type BaseProvider struct {
	Firmeve *Firmeve `inject:"firmeve"`
}

type Firmeve struct {
	container.Container
	providers map[string]Provider
	booted    bool
	mode      uint8
}

type registerOption struct {
	registerForce bool
}

type option struct {
	mode uint8
}

func WithMode(mode uint8) support.Option {
	return func(object support.Object) {
		object.(*option).mode = mode
	}
}

// Create a new firmeve container
func New(options ...support.Option) *Firmeve {

	option := support.ApplyOption(newOption(), options...).(*option)

	firmeve := &Firmeve{
		Container: container.New(),
		providers: make(map[string]Provider),
		booted:    false,
		mode:      option.mode,
	}

	// binding self
	firmeve.Bind("firmeve", firmeve)

	return firmeve
}

// binding unique firmeve instance
//func BindingInstance(firmeve *Firmeve) {
//	if instance != nil {
//		return
//	}
//
//	once.Do(func() {
//		instance = firmeve
//	})
//}

// A singleton firmeve expose func
//func Instance() *Firmeve {
//	return instance
//}
//
//func F(params ...interface{}) interface{} {
//	if len(params) > 0 {
//		return Instance().Make(params[0], params[1:]...)
//	}
//
//	return Instance()
//}

// Set running mode
func (f *Firmeve) SetMode(mode uint8) *Firmeve {
	f.mode = mode
	return f
}

// Get running mode
func (f *Firmeve) Mode() uint8 {
	return f.mode
}

// Check is development mode
func (f *Firmeve) IsDevelopment() bool {
	return f.mode == ModeDevelopment
}

// Check is production mode
func (f *Firmeve) IsProduction() bool {
	return f.mode == ModeProduction
}

// Check is testing mode
func (f *Firmeve) IsTesting() bool {
	return f.mode == ModeTesting
}

// Start all service providers at once
func (f *Firmeve) Boot() {
	if f.booted {
		return
	}

	for _, provider := range f.providers {
		provider.Boot()
	}

	f.booted = true
}

// Compatible method make method alias
func (f *Firmeve) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return f.Make(abstract, params...)
}

// Register force param
func WithRegisterForce() support.Option {
	return func(object support.Object) {
		object.(*registerOption).registerForce = true
	}
}

// Register a service provider
func (f *Firmeve) Register(provider Provider, options ...support.Option) {
	name := provider.Name()
	// Parameter analysis
	option := support.ApplyOption(newRegisterOption(), options...).(*registerOption)

	if f.HasProvider(name) && !option.registerForce {
		return
	}

	f.registerProvider(name, provider)

	provider.Register()

	if f.booted {
		provider.Boot()
	}
}

// Add a service provider to providers map
func (f *Firmeve) registerProvider(name string, provider Provider) {
	var mutex sync.Mutex
	mutex.Lock()
	f.providers[name] = provider
	mutex.Unlock()
}

// Determine if the provider exists
func (f *Firmeve) HasProvider(name string) bool {
	if _, ok := f.providers[name]; ok {
		return ok
	}

	return false
}

// Get a if the provider exists
// If not found then panic
func (f *Firmeve) GetProvider(name string) Provider {
	if !f.HasProvider(name) {
		panic(fmt.Errorf("the %s service provider not exists", name))
	}

	return f.providers[name]
}

// ---------------------------- option ------------------------

func newRegisterOption() *registerOption {
	return &registerOption{registerForce: false}
}

func newOption() *option {
	return &option{mode: ModeProduction}
}
