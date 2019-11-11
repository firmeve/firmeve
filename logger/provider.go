package logging

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

var (
	Logger = Default()
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id      int
}

func (p *Provider) Register() {
	//@todo 这里需要引入config
	Logger = Default()
	p.Firmeve.Bind(`logger`, Logger, container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve.Instance().Register(`logger`, firmeve.Instance().Resolve(new(Provider)).(*Provider))
}
