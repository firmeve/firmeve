package http

import (
	"github.com/firmeve/firmeve/container"
	"github.com/gorilla/mux"
	net_http "net/http"
	"sync"
)

type ServiceProvider struct {
	Firmeve *container.Firmeve `inject:"firmeve"`
}

type Http struct {
}

type Router struct {
	router *mux.Router
}

var (
	http *Http
	once sync.Once
)

//func init() {
//	firmeve := container.GetFirmeve()
//	firmeve.Register(`http`, firmeve.GetContainer().Resolve(new(ServiceProvider)).(*ServiceProvider))
//}

func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

//func (r *Router) ServeHttp(w net_http.ResponseWriter, req *net_http.Request)  {
//	if !r.skipClean {
//		path := req.URL.Path
//		if r.useEncodedPath {
//			path = req.URL.EscapedPath()
//		}
//		// Clean path to canonical form and redirect.
//		if p := cleanPath(path); p != path {
//
//			// Added 3 lines (Philip Schlump) - It was dropping the query string and #whatever from query.
//			// This matches with fix in go 1.2 r.c. 4 for same problem.  Go Issue:
//			// http://code.google.com/p/go/issues/detail?id=5252
//			url := *req.URL
//			url.Path = p
//			p = url.String()
//
//			w.Header().Set("Location", p)
//			w.WriteHeader(http.StatusMovedPermanently)
//			return
//		}
//	}
//	var match RouteMatch
//	var handler http.Handler
//	if r.Match(req, &match) {
//		handler = match.Handler
//		req = setVars(req, match.Vars)
//		req = setCurrentRoute(req, match.Route)
//	}
//
//	if handler == nil && match.MatchErr == ErrMethodMismatch {
//		handler = methodNotAllowedHandler()
//	}
//
//	if handler == nil {
//		handler = http.NotFoundHandler()
//	}
//
//	handler.ServeHTTP(w, req)
//}

//func (r *Router) ServeHttp(w net_http.ResponseWriter, req *net_http.Request) {
//	ctxFunc(NewContext(w, r))
//}

func (r *Router) Get(pattern string, ctxFunc func(ctx *Context)) {
	r.router.Methods(`GET`).HandlerFunc(func(w net_http.ResponseWriter, req *net_http.Request) {
		ctxFunc(NewContext(w, req))
	})
}

// --------------------------------- http server -------------------------------------------

//func NewHttp(config *config.Config) *Http {
//	if http != nil {
//		return http
//	}
//
//	once.Do(func() {
//		http = &Http{
//			Server: gin.New(),
//		}
//	})
//
//	return http
//}
//
//func (h *Http) Run() {
//	//h.Server.Run(h.config.Item("server").GetString("http.host"))
//	h.Server.Run(container.GetFirmeve().GetContainer().Get(`config`).(*config.Config).Item("server").GetString("http.host"))
//}

//func (h *Http) Route() {
//	//h.Server.GET(`/test`, func(context *gin.Context) {
//	//	//fmt.Printf("%#v\n",server.NewContext(context))
//	//	context.String(200, `test`)
//	//})
//}
//
//func (sp *ServiceProvider) Register() {
//	// Register Http Server
//	sp.Firmeve.GetContainer().Bind(`http.server`, NewHttp)
//}
//
//func (sp *ServiceProvider) Boot() {
//	// Register Router
//	sp.Firmeve.GetContainer().Get(`http.server`).(*Http).Route()
//
//	// Use Http middleware
//
//	// Run server
//}
