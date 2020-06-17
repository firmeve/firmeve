package cron

import "github.com/firmeve/firmeve/kernel"

type Provider struct {
	kernel.BaseProvider
}

func (p Provider) Name() string {
	return `cron`
}

func (p *Provider) Register() {
	p.Application.Bind(`cron`, New())
}

func (p Provider) Boot() {

}
