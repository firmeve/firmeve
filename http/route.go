package http

import "github.com/firmeve/firmeve/kernel/contract"

type (
	Route struct {
		path     string
		name     string
		handlers []contract.ContextHandler
		handler  contract.ContextHandler
	}
)

func (r *Route) Name(name string) contract.HttpRoute {
	r.name = name
	return r
}

func (r *Route) Use(handlers ...contract.ContextHandler) contract.HttpRoute {
	r.handlers = append(r.handlers, handlers...)
	return r
}

func (r *Route) Handlers() []contract.ContextHandler {
	return append(r.handlers, r.handler)
}

func newRoute(path string, handler contract.ContextHandler) contract.HttpRoute {
	return &Route{
		path:     path,
		handler:  handler,
		handlers: make([]contract.ContextHandler, 0),
	}
}
