package database

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
)

type Provider struct {
	firmeve.BaseFirmeve
}

func (p *Provider) Name() string {
	return `db`
}

func (p *Provider) Register() {
	DB := New(p.Firmeve.Get(`config`).(config.Configurator).Item(`database`))
	p.Firmeve.Bind(`db`, DB)
	p.Firmeve.Bind(`db.connection`, DB.ConnectionDefault())
}

func (p *Provider) Boot() {
}
