package scheduler

import "github.com/firmeve/firmeve/kernel"

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `scheduler`
}

func (p *Provider) Register() {
	config := new(Configuration)
	p.Config.Bind(`scheduler`, config)

	p.Bind(`scheduler`, New(config))
}

func (p *Provider) Boot() {
}
