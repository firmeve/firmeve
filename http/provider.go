package http

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/gorilla/sessions"
)

type Provider struct {
	kernel.BaseProvider
	Config *config2.Config
}

func (p *Provider) Name() string {
	return `http`
}

func (p *Provider) Register() {
	frameworkConfig := p.Config.Item("framework")

	p.Firmeve.Bind(`http.router`, New(p.Firmeve), container.WithShare(true))

	// session
	p.Firmeve.Bind(`http.session.store`, sessions.NewCookieStore(
		[]byte(frameworkConfig.GetString("key")),
	))
}

func (p *Provider) Boot() {

}
