package http

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `http`
}

func (p *Provider) Register() {
	p.Application.RegisterPool(`http`, func(application contract.Application) interface{} {
		return &Http{
			application: application,
			params:      make([]httprouter.Param, 3),
		}
	})

	p.Application.Bind(`http.router`, New(p.Application), container.WithShare(true))

	// session
	p.Application.Bind(`http.session.store`, sessions.NewCookieStore(
		[]byte(p.Config.GetString(`framework.key`)),
	))
}

func (p *Provider) Boot() {

}
