package logging

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `logger`
}

func (p *Provider) Register() {
	config := p.Firmeve.Get(`config`).(*config2.Config).Item(`logging`)
	p.Bind(`logger`, New(config), container.WithShare(true))
}

func (p *Provider) Boot() {

}
