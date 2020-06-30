package main

import (
	"errors"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
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
			c.Error(401, errors.New(`未认证`))
		})
	}
}

func main() {
	firmeve.RunDefault(firmeve.WithConfigPath("../config.yaml"), firmeve.WithProviders(
		[]contract.Provider{
			new(App),
		},
	))
}
