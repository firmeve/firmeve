package http

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"net/http"
	"strings"
)

type Group struct {
	prefix         string
	beforeHandlers []contract.ContextHandler
	afterHandlers  []contract.ContextHandler
	router         *Router
}

func (g *Group) Prefix(prefix string) *Group {
	g.prefix = prefix
	return g
}

func (g *Group) After(handlers ...contract.ContextHandler) *Group {
	g.afterHandlers = append(g.afterHandlers, handlers...)
	return g
}

func (g *Group) Before(handlers ...contract.ContextHandler) *Group {
	g.beforeHandlers = append(g.beforeHandlers, handlers...)
	return g
}

func (g *Group) GET(path string, handler contract.ContextHandler) *Route {
	return g.createRoute(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler contract.ContextHandler) *Route {
	return g.createRoute(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler contract.ContextHandler) *Route {
	return g.createRoute(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler contract.ContextHandler) *Route {
	return g.createRoute(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler contract.ContextHandler) *Route {
	return g.createRoute(http.MethodDelete, path, handler)
}

func (g *Group) OPTIONS(path string, handler contract.ContextHandler) *Route {
	return g.createRoute(http.MethodOptions, path, handler)
}

func (g *Group) Group(prefix string) *Group {
	return newGroup(g.router).Prefix(strings.Join([]string{g.prefix, prefix}, ``)).After(g.afterHandlers...).Before(g.beforeHandlers...)
}

func (g *Group) createRoute(method string, path string, handler contract.ContextHandler) *Route {
	path = strings.Join([]string{g.prefix, path}, ``)

	return g.router.createRoute(method, path, handler).Before(g.beforeHandlers...).After(g.afterHandlers...)
}

func newGroup(router *Router) *Group {
	return &Group{
		router:         router,
		beforeHandlers: make([]contract.ContextHandler, 0),
		afterHandlers:  make([]contract.ContextHandler, 0),
	}
}
