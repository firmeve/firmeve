package bootstrap

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cache"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/database"
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/logger"
	"github.com/firmeve/firmeve/support"
	"github.com/kataras/iris/core/errors"
	"os"
)

type (
	Bootstrap struct {
		Firmeve    *firmeve.Firmeve
		configPath string
	}

	option struct {
		path string
	}

	RegisterFunc func(b *Bootstrap)
)

var (
	configPathError = errors.New(`config path error`)
)

func WithConfigPath(path string) support.Option {
	return func(object support.Object) {
		object.(*option).path = path
	}
}

func New(firmeve2 *firmeve.Firmeve, options ...support.Option) *Bootstrap {
	option := support.ApplyOption(&option{}, options...).(*option)

	if option.path == `` {
		option.path = os.Getenv("FIRMEVE_CONFIG_PATH")
		if option.path == `` {
			panic(configPathError)
		}
	}

	bootstrap := &Bootstrap{
		Firmeve:    firmeve2,
		configPath: option.path,
	}
	bootstrap.configure()
	bootstrap.registerBaseProvider()
	return bootstrap
}

func (b *Bootstrap) RegisterDefault() *Bootstrap {
	return b.Register([]firmeve.Provider{
		b.Firmeve.Make(new(cache.Provider)).(firmeve.Provider),
		b.Firmeve.Make(new(database.Provider)).(firmeve.Provider),
		b.Firmeve.Make(new(http.Provider)).(firmeve.Provider),
	}, firmeve.WithRegisterForce())
}

func (b *Bootstrap) Register(providers []firmeve.Provider, options ...support.Option) *Bootstrap {
	for _, provider := range providers {
		b.Firmeve.Register(provider, options...)
	}

	return b
}

func (b *Bootstrap) Boot() {
	b.Firmeve.Boot()
}

func (b *Bootstrap) FastBootFull() {
	b.RegisterDefault()
	b.Boot()
}

func (b *Bootstrap) FastBootFullWithProviders(providers []firmeve.Provider, options ...support.Option) {
	b.RegisterDefault()

	newProviders := make([]firmeve.Provider, 0)
	for _, provider := range providers {
		newProviders = append(newProviders, b.Firmeve.Make(provider).(firmeve.Provider))
	}

	b.Register(newProviders, options...)

	b.Boot()
}

func (b *Bootstrap) configure() {
	b.Firmeve.Bind(`config`, config2.New(b.configPath), container.WithShare(true))
}

func (b *Bootstrap) registerBaseProvider() {
	b.Register([]firmeve.Provider{
		b.Firmeve.Make(new(event.Provider)).(firmeve.Provider),
		b.Firmeve.Make(new(logging.Provider)).(firmeve.Provider),
	})
}
