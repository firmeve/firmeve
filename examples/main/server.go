package main

import (
	"os"
	path2 "github.com/firmeve/firmeve/support/path"
	"github.com/firmeve/firmeve/bootstrap"
	cmd2 "github.com/firmeve/firmeve/http/cmd"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cmd"
	"github.com/firmeve/firmeve/http"
	_ "github.com/takama/daemon"
)

type App struct {
	firmeve.BaseProvider
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
		v1.GET(`/ping`, func(c *http.Context) {
			c.Data(map[string]string{
				"message": "pong",
			})
			c.Next()
		})
	}
}

func main() {
	// Base bootstrap
	app := firmeve.New()
	bootstrap2 := bootstrap.New(app, bootstrap.WithConfigPath(path2.RunRelative(`../../testdata/config`)))
	bootstrap2.FastBootFullWithProviders(
		[]firmeve.Provider{new(App)},
	)

	// Command
	root := cmd.Root()
	root.AddCommand(cmd2.NewServer(bootstrap2).Cmd())
	root.SetArgs(os.Args[1:])
	root.Execute()
}
