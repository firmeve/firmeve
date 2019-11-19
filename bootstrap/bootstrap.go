package bootstrap

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cache"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/database"
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/http"
	logging "github.com/firmeve/firmeve/logger"
)

type Bootstrap struct {
	configPath string
	Firmeve    *firmeve.Firmeve
}

func New(firmeve2 *firmeve.Firmeve, configPath string) *Bootstrap {
	//binding current unique firmeve
	//firmeve.BindingInstance(firmeve2)

	bootstrap := &Bootstrap{
		configPath: configPath,
		Firmeve:    firmeve2,
	}
	bootstrap.configure(configPath)
	bootstrap.registerBaseProvider()
	return bootstrap
}

func (b *Bootstrap) RegisterDefault() *Bootstrap {
	return b.Register([]firmeve.Provider{
		b.Firmeve.Make(new(cache.Provider)).(firmeve.Provider),
		b.Firmeve.Make(new(database.Provider)).(firmeve.Provider),
		b.Firmeve.Make(new(http.Provider)).(firmeve.Provider),
	}...)
}

func (b *Bootstrap) Register(providers ...firmeve.Provider) *Bootstrap {
	for _, provider := range providers {
		b.Firmeve.Register(provider, firmeve.WithRegisterForce())
	}

	return b
}

func (b *Bootstrap) Boot() {
	b.Firmeve.Boot()
}

func (b *Bootstrap) configure(path string) {
	b.Firmeve.Bind(`config`, config2.New(path), container.WithShare(true))
}

func (b *Bootstrap) registerBaseProvider() {
	b.Register([]firmeve.Provider{
		b.Firmeve.Make(new(event.Provider)).(firmeve.Provider),
		b.Firmeve.Make(new(logging.Provider)).(firmeve.Provider),
	}...)
}
