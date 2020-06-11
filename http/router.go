package http

import (
	context2 "context"
	"github.com/firmeve/firmeve/context"
	"github.com/firmeve/firmeve/kernel/contract"
	http2 "github.com/firmeve/firmeve/support/http"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	time2 "time"
)

type Router struct {
	Application contract.Application
	router      *httprouter.Router
	routes      map[string]contract.HttpRoute
	routeKeys   []string
	event       contract.Event
	logger      contract.Loggable
}

func New(app contract.Application) contract.HttpRouter {
	return &Router{
		Application: app,
		router:      httprouter.New(),
		routes:      make(map[string]contract.HttpRoute, 0),
		routeKeys:   make([]string, 0),
		event:       app.Resolve(`event`).(contract.Event),
		logger:      app.Resolve(`logger`).(contract.Loggable),
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
		context.NewContext(r.Application, NewHttp(r.Application, req, w.(contract.HttpWrapResponseWriter)), handler).Next()
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
		// base http
		currentHttp := NewHttp(r.Application, req, w.(contract.HttpWrapResponseWriter))
		currentHttp.SetParams(params)
		currentHttp.SetRoute(r.routes[key])

		// context create
		ctx := context.NewContext(r.Application, currentHttp, r.routes[key].Handlers()...)

		// router match dispatch
		r.event.Dispatch(`router.match`, map[string]interface{}{
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
	// wrap record statusCode
	wrap := NewWrapResponseWriter(w)

	// dispatch router
	r.event.Dispatch(`http.request`, map[string]interface{}{
		`request`:  req.Clone(context2.Background()),
		`response`: wrap,
	})

	// record start time
	startTime := time2.Now()

	r.router.ServeHTTP(wrap, req)

	// request log
	r.logger.Debug(``,
		`Method`, req.Method,
		`StatusCode`, wrap.StatusCode(),
		`URI`, req.RequestURI,
		`IPAddress`, http2.ClientIP(req),
		`Agent`, req.Header.Get(`user-agent`),
		`ExecuteTime`, time2.Now().Sub(startTime),
	)
}
