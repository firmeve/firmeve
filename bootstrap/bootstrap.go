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
)

type bootstrap struct {
	configPath string
	firmeve    *firmeve.Firmeve
}

func New(firmeve2 *firmeve.Firmeve, configPath string) *bootstrap {
	//binding current unique firmeve
	//firmeve.BindingInstance(firmeve2)

	return &bootstrap{
		configPath: configPath,
		firmeve:    firmeve2,
	}
}

func (b *bootstrap) BootstrapBase() *bootstrap {
	b.configure(b.configPath)
	b.registerBaseProvider()
	return b
}

func (b *bootstrap) RegisterDefault() *bootstrap {
	return b.Register([]firmeve.Provider{
		b.firmeve.Make(new(cache.Provider)).(firmeve.Provider),
		b.firmeve.Make(new(database.Provider)).(firmeve.Provider),
		b.firmeve.Make(new(http.Provider)).(firmeve.Provider),
	}...)
}

func (b *bootstrap) Register(providers ...firmeve.Provider) *bootstrap {
	for _, provider := range providers {
		b.firmeve.Register(provider, firmeve.WithRegisterForce())
	}

	return b
}

func (b *bootstrap) Boot() *bootstrap {
	b.firmeve.Boot()

	return b
}

func (b *bootstrap) QuickBoot() *bootstrap {
	b.BootstrapBase().RegisterDefault().Boot()
	return b
}

func (b *bootstrap) configure(path string) {
	b.firmeve.Bind(`config`, config2.New(path), container.WithShare(true))
}

func (b *bootstrap) registerBaseProvider() {
	b.Register([]firmeve.Provider{
		b.firmeve.Make(new(event.Provider)).(firmeve.Provider),
		b.firmeve.Make(new(logging.Provider)).(firmeve.Provider),
	}...)
}
