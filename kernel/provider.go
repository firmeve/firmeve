package kernel

type (
	IProvider interface {
		Name() string
		Register()
		Boot()
	}

	//registerOption struct {
	//	registerForce bool
	//}

	BaseProvider struct {
		Firmeve IApplication `inject:"firmeve"`
	}
)

// Register force param
//func WithRegisterForce() support.Option {
//	return func(object support.Object) {
//		object.(*registerOption).registerForce = true
//	}
//}

//
//// Register a service provider
//func (p *Provider) Register(provider IProvider, options ...support.Option) {
//	name := provider.Name()
//	// Parameter analysis
//	option := support.ApplyOption(newRegisterOption(), options...).(*registerOption)
//
//	if p.HasProvider(name) && !option.registerForce {
//		return
//	}
//
//	p.registerProvider(name, provider)
//
//	provider.Register()
//
//	if p.booted {
//		provider.Boot()
//	}
//}
//
//// Add a service provider to providers map
//func (p *Provider) registerProvider(name string, provider IProvider) {
//	var mutex sync.Mutex
//	mutex.Lock()
//	p.providers[name] = provider
//	mutex.Unlock()
//}
//
//// Determine if the provider exists
//func (p *Provider) HasProvider(name string) bool {
//	if _, ok := p.providers[name]; ok {
//		return ok
//	}
//
//	return false
//}
//
//// Get a if the provider exists
//// If not found then panic
//func (p *Provider) GetProvider(name string) IProvider {
//	if !p.HasProvider(name) {
//		panic(fmt.Errorf("the %s service provider not exists", name))
//	}
//
//	return p.providers[name]
//}
//
////// ---------------------------- option ------------------------
//
//func newRegisterOption() *registerOption {
//	return &registerOption{registerForce: false}
//}
