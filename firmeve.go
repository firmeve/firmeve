package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/cmd"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/http"
	"github.com/spf13/cobra"
	"os"

	//"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/logger"
	"github.com/firmeve/firmeve/support"
	"sync"
	_ "sync"
	//"github.com/firmeve/firmeve/container"
	//"github.com/firmeve/firmeve/support"
)

type (
	Firmeve struct {
		container.Container
		providers map[string]kernel.IProvider
		booted    bool
		mode      uint8
	}

	option struct {
		providers []kernel.IProvider
	}
)

var (
	registerMutex sync.Mutex
)

//func WithMode(mode uint8) support.Option {
//	return func(object support.Object) {
//		object.(*option).mode = mode
//	}
//}
//
//func WithConfigPath(path string) support.Option {
//	return func(object support.Object) {
//		object.(*option).configPath = path
//	}
//}

func WithProviders(providers []kernel.IProvider) support.Option {
	return func(object support.Object) {
		object.(*option).providers = providers
	}
}

func New(mode uint8, configPath string, options ...support.Option) kernel.IApplication {
	return startNewApplication(mode,configPath,parseOption(options))
}

func Default(mode uint8, configPath string, options ...support.Option) kernel.IApplication {
	defaultProviders := []kernel.IProvider{
		new(http.Provider),
	}

	option := parseOption(options)
	option.providers = append(defaultProviders,option.providers...)

	return startNewApplication(mode, configPath,option)
}

func startNewApplication(mode uint8, configPath string, option *option) kernel.IApplication {
	f := &Firmeve{
		Container: container.New(),
		providers: make(map[string]kernel.IProvider, 0),
		booted:    false,
		mode:      mode,
	}

	f.Bind("firmeve", f)

	f.Bind(`config`, config.New(configPath), container.WithShare(true))

	f.registerBaseProvider(configPath)

	if len(option.providers) != 0 {
		f.RegisterMultiple(option.providers, false)
	}

	f.Boot()

	return f
}

func parseOption(options []support.Option) *option {
	return support.ApplyOption(&option{
		providers: make([]kernel.IProvider, 0),
	}, options...).(*option)
}

func (f *Firmeve) Run() {
	root := f.Get("command").(*cobra.Command)
	root.SetArgs(os.Args[1:])
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}

func (f *Firmeve) SetMode(mode uint8) {
	f.mode = mode
}

func (f *Firmeve) Mode() uint8 {
	return f.mode
}

func (f *Firmeve) IsDevelopment() bool {
	return f.mode == kernel.ModeDevelopment
}

func (f *Firmeve) IsProduction() bool {
	return f.mode == kernel.ModeProduction
}

func (f *Firmeve) IsTesting() bool {
	return f.mode == kernel.ModeTesting
}

func (f *Firmeve) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return f.Make(abstract, params...)
}

func (f *Firmeve) Boot() {
	if f.booted {
		return
	}

	for i := range f.providers {
		f.providers[i].Boot()
	}

	f.booted = true
}

func (f *Firmeve) Register(provider kernel.IProvider, force bool) {
	name := provider.Name()

	if f.HasProvider(name) && !force {
		return
	}

	f.registerProvider(name, provider)

	provider.Register()

	if f.booted {
		provider.Boot()
	}
}

func (f *Firmeve) RegisterMultiple(providers []kernel.IProvider, force bool) {
	for i := range providers {
		f.Register(f.Make(providers[i]).(kernel.IProvider), force)
	}
}

func (f *Firmeve) HasProvider(name string) bool {
	if _, ok := f.providers[name]; ok {
		return ok
	}

	return false
}

func (f *Firmeve) GetProvider(name string) kernel.IProvider {
	if !f.HasProvider(name) {
		panic(fmt.Errorf("the %s service provider not exists", name))
	}

	return f.providers[name]
}

func (f *Firmeve) Reset() {
	f.providers = make(map[string]kernel.IProvider, 0)
	f.Container.Flush()
	f.booted = false
}

// Add a service provider to providers map
func (f *Firmeve) registerProvider(name string, provider kernel.IProvider) {
	registerMutex.Lock()
	f.providers[name] = f.Make(provider).(kernel.IProvider)
	registerMutex.Unlock()
}

func (f *Firmeve) registerBaseProvider(configPath string) {

	f.RegisterMultiple([]kernel.IProvider{
		new(cmd.Provider),
		new(event.Provider),
		new(logging.Provider),
	}, false)
	return
}

//func New(options ...support.Option) {
//	option := support.ApplyOption(&option{}, options...).(*option)
//	app := New2(option.mode, "")
//
//	//f.Bind(`firmeve`, app)
//
//	//loading config
//	//configProvider := f.Resolve(new(config.Provider)).(*config.Provider)
//	//configProvider.ConfigPath = option.configPath
//
//	//f.Register(new(event.Provider), false)
//
//	//baseProvider
//	//f.Register()
//	//f.Bind("config").configPath string
//	////providers
//	//f.Register(new(app2.Provider), false)
//	//f := kernel.Application{
//	//	Container: container.New(),
//	//	//provider:  provider.Provider{},
//	//	providers: make(map[string]provider.IProvider, 0),
//	//	booted:    false,
//	//	//mode:      option.mode,
//	//}
//	//f.provider.Register(new(f.Provider))
//	fmt.Println(app)
//}

//func (f *Firmeve) Register() {
//
//}

//
//// Create a new firmeve container
//func New(options ...support.Option) *Firmeve {
//
//	option := support.ApplyOption(newOption(), options...).(*option)
//
//	firmeve := &Firmeve{
//		Container: container.New(),
//		providers: make(map[string]Provider),
//		booted:    false,
//		mode:      option.mode,
//	}
//
//	// binding self
//	firmeve.Bind("firmeve", firmeve)
//
//	return firmeve
//}
//
//// binding unique firmeve instance
////func BindingInstance(firmeve *Firmeve) {
////	if instance != nil {
////		return
////	}
////
////	once.Do(func() {
////		instance = firmeve
////	})
////}
//
//// A singleton firmeve expose func
////func Instance() *Firmeve {
////	return instance
////}
////
////func F(params ...interface{}) interface{} {
////	if len(params) > 0 {
////		return Instance().Make(params[0], params[1:]...)
////	}
////
////	return Instance()
////}
//
//// Set running mode
//func (f *Firmeve) SetMode(mode uint8) *Firmeve {
//	f.mode = mode
//	return f
//}
//
//// Get running mode
//func (f *Firmeve) Mode() uint8 {
//	return f.mode
//}
//
//// Check is development mode
//func (f *Firmeve) IsDevelopment() bool {
//	return f.mode == ModeDevelopment
//}
//
//// Check is production mode
//func (f *Firmeve) IsProduction() bool {
//	return f.mode == ModeProduction
//}
//
//// Check is testing mode
//func (f *Firmeve) IsTesting() bool {
//	return f.mode == ModeTesting
//}
//
//// Start all service providers at once
//func (f *Firmeve) Boot() {
//	if f.booted {
//		return
//	}
//
//	for _, provider := range f.providers {
//		provider.Boot()
//	}
//
//	f.booted = true
//}
//
//// Compatible method make method alias
//func (f *Firmeve) Resolve(abstract interface{}, params ...interface{}) interface{} {
//	return f.Make(abstract, params...)
//}
//
//// Register force param
//func WithRegisterForce() support.Option {
//	return func(object support.Object) {
//		object.(*registerOption).registerForce = true
//	}
//}
//
//// Register a service provider
//func (f *Firmeve) Register(provider Provider, options ...support.Option) {
//	name := provider.Name()
//	// Parameter analysis
//	option := support.ApplyOption(newRegisterOption(), options...).(*registerOption)
//
//	if f.HasProvider(name) && !option.registerForce {
//		return
//	}
//
//	f.registerProvider(name, provider)
//
//	provider.Register()
//
//	if f.booted {
//		provider.Boot()
//	}
//}
//
//// Add a service provider to providers map
//func (f *Firmeve) registerProvider(name string, provider Provider) {
//	var mutex sync.Mutex
//	mutex.Lock()
//	f.providers[name] = provider
//	mutex.Unlock()
//}
//
//// Determine if the provider exists
//func (f *Firmeve) HasProvider(name string) bool {
//	if _, ok := f.providers[name]; ok {
//		return ok
//	}
//
//	return false
//}
//
//// Get a if the provider exists
//// If not found then panic
//func (f *Firmeve) GetProvider(name string) Provider {
//	if !f.HasProvider(name) {
//		panic(fmt.Errorf("the %s service provider not exists", name))
//	}
//
//	return f.providers[name]
//}
//
//// ---------------------------- option ------------------------
//
//func newRegisterOption() *registerOption {
//	return &registerOption{registerForce: false}
//}
//
//func newOption() *option {
//	return &option{mode: ModeProduction}
//}
