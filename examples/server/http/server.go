package main

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
	path2 "github.com/firmeve/firmeve/support/path"
)

type App struct {
	kernel.BaseProvider
}

func (a *App) Name() string {
	return `app`
}

func (a *App) Register() {
}

func (a *App) Boot() {
	a.bindingRoutes()
}

func (a *App) bindingRoutes() {
	router := a.Firmeve.Get(`http.router`).(*http.Router)
	v1 := router.Group("/api/v1")
	{
		v1.GET(`/ping`, func(c contract.Context) {
			c.RenderWith(200, render.JSON, map[string]string{
				"message": "pong",
			})
			c.Next()
		})
	}
}

func main() {
	app := firmeve.Default(contract.ModeDevelopment, path2.RunRelative(`../../../testdata/config`),
		firmeve.WithProviders(
			[]contract.Provider{new(App)},
		))

	app.Run()
}
