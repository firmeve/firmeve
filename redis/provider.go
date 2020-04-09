package redis

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (Provider) Name() string {
	return `redis`
}

func (p *Provider) Register() {
	client := New(p.Firmeve.Resolve(`config`).(*config2.Config).Item(`redis`))
	p.Firmeve.Bind(`redis.client`, client)
	p.Firmeve.Bind(`redis.client.connection`, client.Connection(`default`))
	p.Firmeve.Bind(`redis.client.cluster`, client.Cluster(`default`))
}

func (Provider) Boot() {

}
