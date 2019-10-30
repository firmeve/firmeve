package event

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id int
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`event`, NewDispatcher(), container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve := firmeve.NewFirmeve()
	firmeve.Register(`event`, firmeve.Resolve(new(Provider)).(*Provider))
}
