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
	*mux.Router
	//route  *mux.Route
}

type Route struct {
	*mux.Route
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
		Router:mux.NewRouter(),
	}
}

func (r *Router) Get(pattern string, ctxFunc func(ctx *Context)) *mux.Route {
	return r.HandleFunc(pattern, func(w net_http.ResponseWriter, r *net_http.Request) {
		ctxFunc(NewContext(w, r))
	}).Methods(`GET`)
}
func (r *Router) Post(pattern string, ctxFunc func(ctx *Context)) *mux.Route {
	return r.HandleFunc(pattern, func(w net_http.ResponseWriter, r *net_http.Request) {
		ctxFunc(NewContext(w, r))
	}).Methods(`POST`)
}
func (r *Router) Delete(pattern string, ctxFunc func(ctx *Context)) *mux.Route {
	return r.HandleFunc(pattern, func(w net_http.ResponseWriter, r *net_http.Request) {
		ctxFunc(NewContext(w, r))
	}).Methods(`DELETE`)
}
func (r *Router) Put(pattern string, ctxFunc func(ctx *Context)) *mux.Route {
	return r.HandleFunc(pattern, func(w net_http.ResponseWriter, r *net_http.Request) {
		ctxFunc(NewContext(w, r))
	}).Methods(`PUT`)
}

//func (r *Router) Options(pattern string, ctxFunc func(ctx *Context)) *mux.Route {
//	return r.HandleFunc(pattern,func(w net_http.ResponseWriter, r *net_http.Request) {
//		ctxFunc(NewContext(w, r))
//	}).Methods(`PUT`)
//}

func (r *Router) PathPrefix(prefix string) *mux.Route {
	return r.PathPrefix(prefix)
}

func (r *Router) Group(f func(router *Router)) {
	//route := NewRouter().router.NewRoute().Subrouter()
	//route.PathPrefix("/sub")
	//r = r.NewRoute().Subrouter()
	f(r)
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
