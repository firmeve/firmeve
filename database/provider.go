package database

import (
	"github.com/firmeve/firmeve/kernel"
	"gorm.io/gorm"
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
	// 每次返回一个新的连接函数
	p.Bind(`db.connection.func`, func() func() *gorm.DB {
		return func() *gorm.DB {
			return database.ConnectionNewDB(config.Default)
		}
	})
}

func (p *Provider) Boot() {
}
