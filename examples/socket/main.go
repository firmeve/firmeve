package main

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/socket"
)

func main() {
	firmeve.RunWithSupportFunc(
		application,
		firmeve.WithConfigPath("./config.yaml"),
		//firmeve.WithProviders(),
		firmeve.WithCommands([]contract.Command{
			new(socket.Command),
		}),
	)
}

func application(application contract.Application) {
	//router := application.Resolve(`http.router`).(contract.HttpRouter)
	//router.GET("/", func(c contract.Context) {
	//	fmt.Printf("%t", c.Firmeve() == firmeve.Application)
	//	c.RenderWith(200, render.JSON, map[string]string{
	//		"ctx_application":    fmt.Sprintf("%p", c.Firmeve()),
	//		"global_application": fmt.Sprintf("%p", firmeve.Application),
	//	})
	//})
}
