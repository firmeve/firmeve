package cache

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id      int
}

func (p *Provider) Register() {
	//@todo 这里要接入config
	p.Firmeve.Bind(`cache`, Default(), container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve.Instance().Register(`cache`, firmeve.Instance().Resolve(new(Provider)).(*Provider))
}
