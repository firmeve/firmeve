package http

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

var ROUTER *Router

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id      int
}

func (p *Provider) Register() {
	ROUTER = New()
	p.Firmeve.Bind(`http.router`, ROUTER, container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve.Instance().Register(`http`, firmeve.Instance().Resolve(new(Provider)).(*Provider))
}

func Singleton() *Router {
	return firmeve.Instance().Get(`http.router`).(*Router)
}
