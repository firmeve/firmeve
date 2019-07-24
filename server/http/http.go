package http

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

type ServiceProvider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
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
	h.server.Run(strings.Join([]string{h.config.Item("server").GetString("http.host"),h.config.Item("server").GetString("http.port")}, `:`))
}

func (h *Http) Route() {
	h.server.GET(`/test`, func(context *gin.Context) {
		context.String(200, `test`)
	})
}

func (hsp *ServiceProvider) Register() {
	// Register Http Server
	hsp.Firmeve.Bind(`http.server`, NewHttp)

}

func (hsp *ServiceProvider) Boot() {
	// Register Router
	hsp.Firmeve.Get(`http.server`).(*Http).Route()

	// Use Http middleware

	// Run server
}
