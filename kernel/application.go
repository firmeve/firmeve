package kernel

import (
	"fmt"
	"github.com/firmeve/firmeve/container"
	"sync"
)

const (
	Version = "1.0.0"
	ModeDevelopment uint8 = iota
	ModeProduction
	ModeTesting
)

type (
	IApplication interface {
		container.Container
		SetMode(mode uint8)
		Mode() uint8
		IsDevelopment() bool
		IsProduction() bool
		IsTesting() bool
		Resolve(abstract interface{}, params ...interface{}) interface{}
		Boot()
		Register(provider IProvider, force bool)
		RegisterMultiple(providers []IProvider, force bool)
		HasProvider(name string) bool
		GetProvider(name string) IProvider
		Reset()
	}

	Application struct {
		container.Container
		providers map[string]IProvider
		booted    bool
		mode      uint8
	}
)

var (
	registerMutex sync.Mutex
)

func New(mode uint8) IApplication {
	return &Application{
		Container: container.New(),
		providers: make(map[string]IProvider, 0),
		booted:    false,
		mode:      mode,
	}
}

func (a *Application) SetMode(mode uint8) {
	a.mode = mode
}

func (a *Application) Mode() uint8 {
	return a.mode
}

func (a *Application) IsDevelopment() bool {
	return a.mode == ModeDevelopment
}

func (a *Application) IsProduction() bool {
	return a.mode == ModeProduction
}

func (a *Application) IsTesting() bool {
	return a.mode == ModeTesting
}

func (a *Application) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return a.Make(abstract, params...)
}

func (a *Application) Boot() {
	if a.booted {
		return
	}

	for i := range a.providers {
		a.providers[i].Boot()
	}

	a.booted = true
}

func (a *Application) Register(provider IProvider, force bool) {
	name := provider.Name()

	if a.HasProvider(name) && !force {
		return
	}

	a.registerProvider(name, provider)

	provider.Register()

	if a.booted {
		provider.Boot()
	}
}

func (a *Application) RegisterMultiple(providers []IProvider, force bool) {
	for i := range providers {
		a.Register(a.Make(providers[i]).(IProvider), force)
	}
}

func (a *Application) HasProvider(name string) bool {
	if _, ok := a.providers[name]; ok {
		return ok
	}

	return false
}

func (a *Application) GetProvider(name string) IProvider {
	if !a.HasProvider(name) {
		panic(fmt.Errorf("the %s service provider not exists", name))
	}

	return a.providers[name]
}

func (a *Application) Reset() {
	a.providers = make(map[string]IProvider, 0)
	a.Container.Flush()
	a.booted = false
}

// Add a service provider to providers map
func (a *Application) registerProvider(name string, provider IProvider) {
	registerMutex.Lock()
	a.providers[name] = a.Make(provider).(IProvider)
	registerMutex.Unlock()
}