package http

import (
	"github.com/firmeve/firmeve/container"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
	"sync"
)

type ServiceProvider struct {
	Firmeve *container.Firmeve `inject:"firmeve"`
}

type Http struct {
	Server *gin.Engine
}

type Router struct {
	router *chi.Mux
}

var (
	http *Http
	once sync.Once
)
//
//func init() {
//	firmeve := container.GetFirmeve()
//	firmeve.Register(`http`, firmeve.GetContainer().Resolve(new(ServiceProvider)).(*ServiceProvider))
//}
//
//func NewRouter() *Router {
//	return &Router{
//		router: chi.NewRouter(),
//	}
//}
//
////func (r *Router) Get(pattern string,z func(ctx server.Context)) {
////	r.router.Get(`/`,func(w http.ResponseWriter, r *http.Request) {
////		//context := ctx(w http.ResponseWriter, r *http.Request)
////		//context
////		z(context)
////		//w.Write([]byte("hello world"))
////	})
////}
//
//// --------------------------------- http server -------------------------------------------
//
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
//
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
