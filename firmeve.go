package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/support"
	"sync"
)

type Provider interface {
	Register()
	Boot()
}

type Firmeve struct {
	container.Container

	bashPath string
	//container        *Instance
	providers map[string]Provider
	booted    bool
}

type option struct {
	registerForce bool
}

var (
	instance *Firmeve
	once     sync.Once
)

// Create a new firmeve container
func New() *Firmeve {

	//basePath := os.Getenv("FIRMEVE_BASE_PATH")

	//basePath, err := filepath.Abs(basePath)
	//if err != nil {
	//	panic(err)
	//}

	firmeve := &Firmeve{
		bashPath:  "",
		providers: make(map[string]Provider),
		booted:    false,
		Container: container.New(),
	}

	// binding self
	firmeve.Bind("firmeve", firmeve)

	return firmeve
}

// A singleton firmeve expose func
func Instance() *Firmeve {
	if instance != nil {
		return instance
	}

	once.Do(func() {
		instance = New()
	})

	return instance
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

// Register force param
func WithRegisterForce() support.Option {
	return func(object support.Object) {
		object.(*option).registerForce = true
	}
}

// Register a service provider
func (f *Firmeve) Register(name string, provider Provider, options ...support.Option) {
	// Parameter analysis
	option := support.ApplyOption(newOption(), options...).(*option)

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

// Get application base path
//func (f *Firmeve) GetBasePath() string {
//	return f.bashPath
//}

// ---------------------------- option ------------------------

func newOption() *option {
	return &option{registerForce: false}
}
