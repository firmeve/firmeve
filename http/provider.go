package http

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `http`
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`http.router`, New(p.Firmeve), container.WithShare(true))
}

func (p *Provider) Boot() {

}
