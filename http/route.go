package http

import "github.com/firmeve/firmeve/kernel/contract"

type Route struct {
	path           string
	name           string
	beforeHandlers []contract.ContextHandler
	afterHandlers  []contract.ContextHandler
	handler        contract.ContextHandler
}

func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}

func (r *Route) Before(handlers ...contract.ContextHandler) *Route {
	r.beforeHandlers = append(r.beforeHandlers, handlers...)
	return r
}

func (r *Route) After(handlers ...contract.ContextHandler) *Route {
	r.afterHandlers = append(r.afterHandlers, handlers...)
	return r
}

func (r *Route) Handlers() []contract.ContextHandler {
	return append(append(r.beforeHandlers, r.handler), r.afterHandlers...)
}

func newRoute(path string, handler contract.ContextHandler) *Route {
	return &Route{
		path:           path,
		handler:        handler,
		beforeHandlers: make([]contract.ContextHandler, 0),
		afterHandlers:  make([]contract.ContextHandler, 0),
	}
}
