package jwt

import (
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p Provider) Name() string {
	return `jwt`
}

func (p *Provider) Register() {
	config := new(Configuration)
	p.Config.Bind(`jwt`, config)
	config.Secret = p.Config.GetString(`framework.key`)
	// binding jwt
	p.Bind(`jwt`, New(
		config,
		NewMemoryStore(),
	))
}

func (p Provider) Boot() {
}
