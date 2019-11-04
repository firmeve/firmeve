package http

import (
	"net/http"
	"strings"
)

type Group struct {
	prefix         string
	beforeHandlers []HandlerFunc
	afterHandlers  []HandlerFunc
	router         *Router
}

func (g *Group) Prefix(prefix string) *Group {
	g.prefix = prefix
	return g
}

func (g *Group) After(handlers ...HandlerFunc) *Group {
	g.afterHandlers = append(g.afterHandlers, handlers...)
	return g
}

func (g *Group) Before(handlers ...HandlerFunc) *Group {
	g.beforeHandlers = append(g.beforeHandlers, handlers...)
	return g
}

func (g *Group) GET(path string, handler HandlerFunc) *Route {
	return g.createRoute(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) *Route {
	return g.createRoute(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc) *Route {
	return g.createRoute(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler HandlerFunc) *Route {
	return g.createRoute(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc) *Route {
	return g.createRoute(http.MethodDelete, path, handler)
}

func (g *Group) OPTIONS(path string, handler HandlerFunc) *Route {
	return g.createRoute(http.MethodOptions, path, handler)
}

func (g *Group) Group(prefix string) *Group {
	return newGroup(g.router).Prefix(strings.Join([]string{g.prefix, prefix}, ``)).After(g.afterHandlers...).Before(g.beforeHandlers...)
}

func (g *Group) createRoute(method string, path string, handler HandlerFunc) *Route {
	path = strings.Join([]string{g.prefix, path}, ``)

	return g.router.createRoute(method, path, handler).Before(g.beforeHandlers...).After(g.afterHandlers...)
}

func newGroup(router *Router) *Group {
	return &Group{
		router:         router,
		beforeHandlers: make([]HandlerFunc, 0),
		afterHandlers:  make([]HandlerFunc, 0),
	}
}
