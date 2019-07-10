package firmeve

import (
	"github.com/firmeve/firmeve/config"
	"github.com/uber-go/fx"
	"os"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
	mutex   sync.Mutex
)

type ServiceProvider interface {
	register()
	boot()
}

type Firmeve struct {
	app            *fx.App
	bindOptions    []fx.Option
	resolveOptions []fx.Option
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindOptions:    make([]fx.Option, 1),
			resolveOptions: make([]fx.Option, 1),
		}
	})

	return firmeve
}

func (f *Firmeve) Bind(constructors ...interface{}) {
	f.bindOptions = append(f.bindOptions, fx.Provide(constructors))
}

func (f *Firmeve) Resolve(funcs ...interface{}) {
	f.resolveOptions = append(f.resolveOptions, fx.Invoke(funcs))
}

// Tmp
func (f *Firmeve) Run() {
	// config
	f.Bind(func() (*config.Config, error) {
		configEnvPath := os.Getenv("FIRMEVE_ENV_PATH")
		if configEnvPath == `` {
			configEnvPath = "./testdata/conf"
		}

		return config.NewConfig(configEnvPath)
	})

	f.app = fx.New(append(f.bindOptions, f.resolveOptions...)...)
}
