package http

import (
	"github.com/firmeve/firmeve"
	"github.com/gin-gonic/gin"
	"sync"
)

type ServiceProvider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
}

type Http struct {
	server *gin.Engine
}



var (
	http *Http
	once sync.Once
)

func NewHttp() *Http {
	if http != nil {
		return http
	}

	once.Do(func() {
		http = &Http{
			server:gin.New(),
		}
	})

	return http
}

func (h *Http) Run() {
	h.server.Run(":28080")
}

func (h *Http) Route() {
	h.server.GET(`/test`, func(context *gin.Context) {
		context.String(200,`test`)
	})
}

func (hsp *ServiceProvider) Register() {
	// Register Http Server
	hsp.Firmeve.Bind(`http.server`, NewHttp())

}

func (hsp *ServiceProvider) Boot() {
	// Register Router
	hsp.Firmeve.Get(`http.server`).(*Http).Route()

	// Use Http middleware

	// Run server
}
