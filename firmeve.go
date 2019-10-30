package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/support"
	"sync"
)

type ServiceProvider interface {
	Register()
	Boot()
}

type Firmeve struct {
	container.Container

	bashPath string
	//container        *Instance
	serviceProviders map[string]ServiceProvider
	booted           bool
}

type firmeveOption struct {
	registerForce bool
}

var (
	firmeve *Firmeve
	firmeveOnce sync.Once
)

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}
	//basePath := os.Getenv("FIRMEVE_BASE_PATH")

	//basePath, err := filepath.Abs(basePath)
	//if err != nil {
	//	panic(err)
	//}

	firmeveOnce.Do(func() {
		firmeve = &Firmeve{
			bashPath:         nil,
			serviceProviders: make(map[string]ServiceProvider),
			booted:           false,
			Container:        container.NewContainer(),
		}

		//firmeve.bingingBaseInstance()
	})

	return firmeve
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
func WithRegisterForce(force bool) support.Option {
	return func(object support.Object) {
		object.(*firmeveOption).registerForce = force
	}
}

// Register a service provider
func (f *Firmeve) Register(name string, provider ServiceProvider, options ...support.Option) {
	// Parameter analysis
	firmeveOption := support.ApplyOption(newFirmeveOption(), options...).(*firmeveOption)

	if f.HasProvider(name) && !firmeveOption.registerForce {
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
	var mutex sync.Mutex
	mutex.Lock()
	f.serviceProviders[name] = provider
	mutex.Unlock()
}

// Binging base instance
//func (f *Firmeve) bingingBaseInstance() {
//	firmeve.Bind(`container`, f, WithBindShare(true))
//	firmeve.Bind(`firmeve`, f, WithBindShare(true))
//	firmeve.Bind(`config`, config.NewConfig(strings.Join([]string{f.bashPath, `config`}, `/`)), WithBindShare(true))
//	firmeve.Bind(`logger`, logging.NewLogger, WithBindShare(true))
//	firmeve.Bind(`dispatcher`, event.NewDispatcher, WithBindShare(true))
//}

// Determine if the provider exists
func (f *Firmeve) HasProvider(name string) bool {
	if _, ok := f.serviceProviders[name]; ok {
		return ok
	}

	return false
}

// Get a if the provider exists
// If not found then panic
func (f *Firmeve) GetProvider(name string) ServiceProvider {
	if !f.HasProvider(name) {
		panic(fmt.Errorf("the %s service provider not exists", name))
	}

	return f.serviceProviders[name]
}

// Get application base path
//func (f *Firmeve) GetBasePath() string {
//	return f.bashPath
//}

// ---------------------------- firmeveOption ------------------------

func newFirmeveOption() *firmeveOption {
	return &firmeveOption{registerForce: false}
}