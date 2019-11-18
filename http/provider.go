package http

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id      int
}

func (p *Provider) Name() string {
	return `http`
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`http.router`, New(), container.WithShare(true))
}

func (p *Provider) Boot() {

}
