package http

import (
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
	h.server.Run(h.config.Item("server").GetString("http.host"))
}

func (h *Http) Route() {
	h.server.GET(`/test`, func(context *gin.Context) {
		context.String(200, `test`)
	})
}

func (sp *ServiceProvider) Register() {
	// Register Http Server
	sp.Firmeve.GetContainer().Bind(`http.server`, NewHttp)
}

func (sp *ServiceProvider) Boot() {
	// Register Router
	sp.Firmeve.GetContainer().Get(`http.server`).(*Http).Route()

	// Use Http middleware

	// Run server
}
