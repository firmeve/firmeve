package kernel

import (
	"fmt"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
)

const (
	version = "1.0.0"
)

type (
	Application struct {
		contract.Container
		providers map[string]contract.Provider
		booted    bool
		mode      uint8
	}
)

var (
	registerMutex sync.Mutex
)

func New(mode uint8) contract.Application {
	return &Application{
		Container: container.New(),
		providers: make(map[string]contract.Provider, 0),
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

func (a *Application) Version() string {
	return version
}

func (a *Application) IsDevelopment() bool {
	return a.mode == contract.ModeDevelopment
}

func (a *Application) IsProduction() bool {
	return a.mode == contract.ModeProduction
}

func (a *Application) IsTesting() bool {
	return a.mode == contract.ModeTesting
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

func (a *Application) Register(provider contract.Provider, force bool) {
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

func (a *Application) RegisterMultiple(providers []contract.Provider, force bool) {
	for i := range providers {
		a.Register(a.Make(providers[i]).(contract.Provider), force)
	}
}

func (a *Application) HasProvider(name string) bool {
	if _, ok := a.providers[name]; ok {
		return ok
	}

	return false
}

func (a *Application) GetProvider(name string) contract.Provider {
	if !a.HasProvider(name) {
		panic(fmt.Errorf("the %s service provider not exists", name))
	}

	return a.providers[name]
}

func (a *Application) Reset() {
	a.providers = make(map[string]contract.Provider, 0)
	a.Container.Flush()
	a.booted = false
}

// Add a service provider to providers map
func (a *Application) registerProvider(name string, provider contract.Provider) {
	registerMutex.Lock()
	a.providers[name] = a.Make(provider).(contract.Provider)
	registerMutex.Unlock()
}
