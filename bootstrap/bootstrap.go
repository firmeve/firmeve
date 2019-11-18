package bootstrap

import (
	"github.com/firmeve/firmeve"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/logger"
)

type bootstrap struct {
	configPath string
	firmeve    *firmeve.Firmeve
}

func New(configPath string, firmeve2 *firmeve.Firmeve) *bootstrap {
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

func (b *bootstrap) configure(path string) {
	b.firmeve.Bind(`config`, config2.New(path), container.WithShare(true))
}

func (b *bootstrap) registerBaseProvider() {
	b.firmeve.Register(b.firmeve.Make(new(logging.Provider)).(firmeve.Provider), firmeve.WithRegisterForce())
}
