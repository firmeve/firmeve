package kernel

import (
	"github.com/firmeve/firmeve/container"
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
		Run()
		Resolve(abstract interface{}, params ...interface{}) interface{}
		Boot()
		Register(provider IProvider, force bool)
		RegisterMultiple(providers []IProvider, force bool)
		HasProvider(name string) bool
		GetProvider(name string) IProvider
		Reset()
	}

	//BaseApplication struct {
	//	container.Container
	//	providers map[string]IProvider
	//	booted    bool
	//	mode      uint8
	//}
)

//
//var (
//	registerMutex sync.Mutex
//)
//
//func New(mode uint8, configPath string) IApplication {
//	app := &Application{
//		Container: container.New(),
//		providers: make(map[string]IProvider, 0),
//		booted:    false,
//		mode:      mode,
//	}
//
//	//app.registerBaseProvider(configPath)
//
//	return app
//}
//
//func (app *Application) SetMode(mode uint8) {
//	app.Get("event").(event.IDispatcher).Dispatch("dss", make(event.InParams, 0))
//	app.mode = mode
//}
//
//func (app *Application) Mode() uint8 {
//	return app.mode
//}
//
//func (app *Application) IsDevelopment() bool {
//	return app.mode == ModeDevelopment
//}
//
//func (app *Application) IsProduction() bool {
//	return app.mode == ModeProduction
//}
//
//func (app *Application) IsTesting() bool {
//	return app.mode == ModeTesting
//}
//
//func (app *Application) Resolve(abstract interface{}, params ...interface{}) interface{} {
//	return app.Make(abstract, params...)
//}
//
//func (app *Application) Boot() {
//	if app.booted {
//		return
//	}
//
//	for i := range app.providers {
//		app.providers[i].Boot()
//	}
//
//	app.booted = true
//}
//
//func (app *Application) Register(provider IProvider, force bool) {
//	name := provider.Name()
//
//	if app.HasProvider(name) && !force {
//		return
//	}
//
//	app.registerProvider(name, provider)
//
//	provider.Register()
//
//	if app.booted {
//		provider.Boot()
//	}
//}
//
//func (app *Application) RegisterMultiple(providers []IProvider, force bool) {
//	for i := range providers {
//		app.Register(app.Make(providers[i]).(IProvider), force)
//	}
//}
//
//func (app *Application) HasProvider(name string) bool {
//	if _, ok := app.providers[name]; ok {
//		return ok
//	}
//
//	return false
//}
//
//func (app *Application) GetProvider(name string) IProvider {
//	if !app.HasProvider(name) {
//		panic(fmt.Errorf("the %s service provider not exists", name))
//	}
//
//	return app.providers[name]
//}
//
//func (app *Application) Reset() {
//	app.providers = make(map[string]IProvider, 0)
//	app.Container.Flush()
//	app.booted = false
//}
//
//// Add a service provider to providers map
//func (app *Application) registerProvider(name string, provider IProvider) {
//	registerMutex.Lock()
//	app.providers[name] = provider
//	registerMutex.Unlock()
//}

//func (app *Application) registerBaseProvider(configPath string) {
//	//configProvider :=
//	_ = app.Resolve(new(provider.Provider)).(*provider.Provider)
//	//configProvider.ConfigPath = option.configPath
//	//app.RegisterMultiple()
//	return
//}
