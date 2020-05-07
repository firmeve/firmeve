package main

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
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
	web := router.Group("")
	{
		web.GET("/", func(c contract.Context) {
			c.RenderWith(200, render.Html, render.Template{
				Name:   "index",
				Data:   "海河小王子",
				Append: nil,
			})
			c.Next()
		})
	}
}

func main() {
	firmeve.RunDefault(firmeve.WithConfigPath("./config.yaml"), firmeve.WithProviders(
		[]contract.Provider{
			new(App),
		},
	))
}
