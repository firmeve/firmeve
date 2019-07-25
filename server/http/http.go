package http

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/gin-gonic/gin"
	"sync"
)

type ServiceProvider struct {
	Firmeve *container.Firmeve `inject:"firmeve"`
}

type Http struct {
	config *config.Config
	server *gin.Engine
}

var (
	http *Http
	once sync.Once
)

func NewHttp(config *config.Config) *Http {
	if http != nil {
		return http
	}

	once.Do(func() {
		http = &Http{
			server: gin.New(),
			config: config,
		}
	})

	return http
}

func (h *Http) Run() {
	fmt.Println(h.config.Item("server").GetString("http.host"))
	h.server.Run(h.config.Item("server").GetString("http.host"))
}

func (h *Http) Route() {
	h.server.GET(`/test`, func(context *gin.Context) {
		context.String(200, `test`)
	})
}

func (hsp *ServiceProvider) Register() {
	// Register Http Server
	hsp.Firmeve.GetContainer().Bind(`http.server`, NewHttp)

}

func (hsp *ServiceProvider) Boot() {
	// Register Router
	hsp.Firmeve.GetContainer().Get(`http.server`).(*Http).Route()

	// Use Http middleware

	// Run server
}
