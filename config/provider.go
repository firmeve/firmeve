package config

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
	path2 "github.com/firmeve/firmeve/support/path"
	"os"
)

type Provider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
	id int
}

func (p *Provider) Register() {
	path := os.Getenv("FIRMEVE_CONFIG_PATH")
	if path == "" {
		path = path2.RunRelative(`../testdata/config`)
	}
	p.Firmeve.Bind(`config`, New(path), container.WithShare(true))
}

func (p *Provider) Boot() {

}

func init() {
	firmeve.Instance().Register(`config`, firmeve.Instance().Resolve(new(Provider)).(*Provider))
}
