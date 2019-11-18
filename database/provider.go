package database

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
}

func (p *Provider) Name() string {
	return `db`
}

func (p *Provider) Register() {
	DB := New(p.Firmeve.Get(`config`).(config.Configurator))
	p.Firmeve.Bind(`db`, DB)
	p.Firmeve.Bind(`db.connection`, DB.ConnectionDefault())
}

func (p *Provider) Boot() {
}
