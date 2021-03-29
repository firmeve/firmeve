package context

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `context`
}

func (p *Provider) Register() {
	p.Application.RegisterPool(`context`, func(application contract.Application) interface{} {
		return &context{
			application: application,
			entries:     make(map[string]*contract.ContextEntity, 0),
			index:       0,
		}
	})
}

func (p *Provider) Boot() {
}
