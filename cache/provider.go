package cache

import (
	"github.com/firmeve/firmeve"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id      int
}

func (p *Provider) Name() string {
	return `cache`
}

func (p *Provider) Register() {
	config := p.Firmeve.Get(`config`).(config2.Configurator).Item(`cache`)
	p.Firmeve.Bind(`cache`, New(&Config{
		Prefix:  config.GetString(`prefix`),
		Current: config.GetString(`default`),
		Repositories: ConfigRepositoryType{
			`redis`: ConfigRepositoryType{
				`host`: config.GetString(`host`),
				`port`: config.GetString(`port`),
				`db`:   config.GetInt(`db`),
			},
		},
	}), container.WithShare(true))
}

func (p *Provider) Boot() {

}
