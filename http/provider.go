package http

import (
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
	p.Application.Bind(`http.router`, New(p.Application), container.WithShare(true))

	// session
	p.Application.Bind(`http.session.store`, sessions.NewCookieStore(
		[]byte(p.Config.GetString(`framework.key`)),
	))
}

func (p *Provider) Boot() {

}
