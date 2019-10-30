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
	p.Firmeve.Bind(`config`, NewConfig(path), container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve := firmeve.NewFirmeve()
	firmeve.Register(`config`, firmeve.Resolve(new(Provider)).(*Provider))
}
