package database

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `db`
}

func (p *Provider) Register() {
	DB := New(p.Firmeve.Get(`config`).(*config.Config).Item(`database`))
	p.Firmeve.Bind(`db`, DB)
	p.Firmeve.Bind(`db.connection`, DB.ConnectionDefault())
}

func (p *Provider) Boot() {
}
