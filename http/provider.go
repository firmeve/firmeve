package http

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id      int
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`http.router`, New(), container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve.Instance().Register(`http`, firmeve.Instance().Resolve(new(Provider)).(*Provider))
}

func Singleton() *Router {
	return firmeve.Instance().Get(`http.router`).(*Router)
}
