package main

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
	"net/http/pprof"
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
	router := a.Firmeve.Get(`http.router`).(contract.HttpRouter)
	v1 := router.Group("/api/v1")
	{
		v1.GET(`/ping`, func(c contract.Context) {
			c.RenderWith(200, render.JSON, map[string]string{
				"message": "pong",
			})
			c.Next()
		})
	}
	debug := router.Group("/debug")
	{
		debug.Handler("GET", "/pprof", pprof.Index)
		debug.Handler("GET", "/cmdline", pprof.Cmdline)
		debug.Handler("GET", "/profile", pprof.Profile)
		debug.Handler("GET", "/symbol", pprof.Symbol)
		debug.Handler("GET", "/trace", pprof.Trace)
	}
}

func main() {
	firmeve.RunDefault(firmeve.WithProviders(
		[]contract.Provider{
			new(App),
		},
	))
}
