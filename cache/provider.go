package cache

import (
	"github.com/firmeve/firmeve/cache/redis"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	redis2 "github.com/firmeve/firmeve/redis"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `cache`
}

func (p *Provider) Register() {
	// register redis
	p.Firmeve.Register(new(redis2.Provider), false)

	config := p.Firmeve.Get(`config`).(*config2.Config).Item(`cache`)
	currentDriver := config.GetString(`default`)
	cache := New(currentDriver)
	//base driver
	cache.Register(`redis`, redis.New(
		p.Resolve(`redis.client`).(*redis2.Redis).Connection(config.GetString(`repositories.`+currentDriver+`.connection`)),
		config.GetString(`prefix`),
	))

	p.Firmeve.Bind(`cache`, cache, container.WithShare(true))

}

func (p *Provider) Boot() {

}
