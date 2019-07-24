package http

import (
	"github.com/firmeve/firmeve"
	"github.com/gin-gonic/gin"
)

type HttpServiceProvider struct {
	Provider *firmeve.FirmeveServiceProvider `inject:"firmeve.provider"`
}

func (hsp *HttpServiceProvider) Register() {
	// Register Http Server
	server := gin.Default()
	hsp.Provider.Firmeve.Bind(`http.server`, server)


}

func (hsp *HttpServiceProvider) Boot() {
	// Register Router

	// Use Http middleware

	// Run server
}
