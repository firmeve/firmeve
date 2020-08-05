package logging

import (
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
	p.Bind(`logger`, New(p.Application), container.WithShare(true))
}

func (p *Provider) Boot() {

}
