package database

import (
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `database`
}

func (p *Provider) Register() {
	var config = new(Configuration)
	p.BindConfig(`database`, config)
	var database = New(config)
	p.Bind(`db`, database)
	// 默认连接
	p.Bind(`db.connection`, database.ConnectionDB(config.Default))
	p.Bind(`db.connection.new`, database.ConnectionNewDB(config.Default))
}

func (p *Provider) Boot() {
}
