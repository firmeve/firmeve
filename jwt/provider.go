package jwt

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
	Config *config2.Config `inject:"config"`
}

func (p Provider) Name() string {
	return `jwt`
}

func (p *Provider) Register() {
	frameworkConfig := p.Config.Item("framework")

	// binding jwt
	p.Bind(`jwt`, New(
		frameworkConfig.GetString("key"),
		p.Config.Item("jwt"),
		NewMemoryStore(),
	))
}

func (p Provider) Boot() {
}
