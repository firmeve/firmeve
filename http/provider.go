package http

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/gorilla/sessions"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `http`
}

func (p *Provider) Register() {
	frameworkConfig := p.Resolve(`config`).(*config2.Config).Item("framework")

	p.Application.Bind(`http.router`, New(p.Application), container.WithShare(true))

	// session
	p.Application.Bind(`http.session.store`, sessions.NewCookieStore(
		[]byte(frameworkConfig.GetString("key")),
	))
}

func (p *Provider) Boot() {

}
