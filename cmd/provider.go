package cmd

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `command`
}

func (p *Provider) Register() {
	p.Firmeve.Bind("command",New(kernel.Version),container.WithShare(true))
}

func (p *Provider) Boot() {
}
