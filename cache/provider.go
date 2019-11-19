package cache

import (
	"github.com/firmeve/firmeve"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	firmeve.BaseFirmeve
}

func (p *Provider) Name() string {
	return `cache`
}

func (p *Provider) Register() {
	config := p.Firmeve.Get(`config`).(*config2.Config).Item(`cache`)
	p.Firmeve.Bind(`cache`, New(config), container.WithShare(true))
}

func (p *Provider) Boot() {

}
