package logging

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `logger`
}

func (p *Provider) Register() {
	config := new(Configuration)
	p.BindConfig(`logging`, config)
	p.Bind(`logger`, New(config, p.Application.Resolve(`event`).(contract.Event)), container.WithShare(true))
}

func (p *Provider) Boot() {

}
