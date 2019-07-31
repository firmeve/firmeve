package http

import (
	"context"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/server"
	"github.com/go-chi/chi"
	net_http "net/http"
	"sync"
)

type ServiceProvider struct {
	Firmeve *container.Firmeve `inject:"firmeve"`
}

type Http struct {
}

type Router struct {
	router *chi.Mux
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
		router: chi.NewRouter(),
	}
}

func (r *Router) Get(pattern string, ctxFunc func(context.Context)) {
	r.router.Get(pattern, func(w net_http.ResponseWriter, r *net_http.Request) {
		ctxFunc(server.NewContext(&Context{request: r, response: w}))
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
