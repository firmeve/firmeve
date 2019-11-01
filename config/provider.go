package config

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
	"os"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id int
}

func (p *Provider) Register() {
	path := os.Getenv("FIRMEVE_CONFIG_PATH")
	if path == "" {
		path = "../testdata/config"
	}
	p.Firmeve.Bind(`config`, New(path), container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve.Instance().Register(`config`, firmeve.Instance().Resolve(new(Provider)).(*Provider))
}
