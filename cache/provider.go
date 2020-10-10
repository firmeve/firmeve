package cache

import (
	"github.com/firmeve/firmeve/cache/redis"
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
	p.Application.Register(new(redis2.Provider), false)

	config := new(Configuration)
	p.BindConfig(`cache`, config)

	cache := New(config)
	redisClient := p.Resolve(`redis`).(*redis.Redis).Client(config.Repositories[`redis`].Connection)
	//config.Default = redis
	cache.Register(config.Default, redis.New(redisClient, config.Prefix))

	p.Application.Bind(`cache`, cache, container.WithShare(true))

}

func (p *Provider) Boot() {

}
