package http

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type Router struct {
	Firmeve   contract.Application
	router    *httprouter.Router
	routes    map[string]contract.HttpRoute
	routeKeys []string
}

func New(firmeve contract.Application) contract.HttpRouter {
	return &Router{
		Firmeve:   firmeve,
		router:    httprouter.New(),
		routes:    make(map[string]contract.HttpRoute, 0),
		routeKeys: make([]string, 0),
	}
}

func (r *Router) GET(path string, handler contract.ContextHandler) contract.HttpRoute {
	return r.createRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler contract.ContextHandler) contract.HttpRoute {
	return r.createRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler contract.ContextHandler) contract.HttpRoute {
	return r.createRoute(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler contract.ContextHandler) contract.HttpRoute {
	return r.createRoute(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler contract.ContextHandler) contract.HttpRoute {
	return r.createRoute(http.MethodDelete, path, handler)
}

func (r *Router) OPTIONS(path string, handler contract.ContextHandler) contract.HttpRoute {
	return r.createRoute(http.MethodOptions, path, handler)
}

// serve static files
func (r *Router) Static(path string, root string) contract.HttpRouter {
	r.router.ServeFiles(strings.Join([]string{path, `/*filepath`}, ``), http.Dir(root))
	return r
}

func (r *Router) NotFound(handler contract.ContextHandler) contract.HttpRouter {
	r.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		kernel.NewContext(r.Firmeve, NewHttp(req, w), handler).Next()
	})

	return r
}

func (r *Router) Handler(method, path string, handler http.HandlerFunc) {
	r.createRoute(method, path, func(c contract.Context) {
		protocol := c.Protocol().(contract.HttpProtocol)
		handler(protocol.ResponseWriter(), protocol.Request())
	})
}

func (r *Router) HttpRouter() *httprouter.Router {
	return r.router
}

//
func (r *Router) Group(prefix string) contract.HttpRouteGroup {
	return newGroup(r).Prefix(prefix)
}

func (r *Router) createRoute(method string, path string, handler contract.ContextHandler) contract.HttpRoute {
	key := r.routeKey(method, path)
	r.routes[key] = newRoute(path, handler)

	r.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		currentHttp := NewHttp(req, w)
		currentHttp.SetParams(params)
		currentHttp.SetRoute(r.routes[key])

		ctx := kernel.NewContext(r.Firmeve, currentHttp, r.routes[key].Handlers()...)

		r.Firmeve.Get(`event`).(contract.Event).Dispatch(`router.match`, map[string]interface{}{
			`context`: ctx,
			`route`:   r.routes[key],
		})

		ctx.Next()
	})

	return r.routes[key]
}

func (r *Router) routeKey(method, path string) string {
	return strings.Join([]string{method, path}, `.`)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
