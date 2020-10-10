package main

import (
	"context"
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
	"net/http/pprof"
	"time"
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
	router := a.Application.Get(`http.router`).(contract.HttpRouter)
	v1 := router.Group("/api/v1")
	{
		v1.Use(http.Recovery)
		v1.GET(`/ping`, func(c contract.Context) {
			c.RenderWith(200, render.JSON, map[string]string{
				"message": "pong",
			})
			c.Next()
		})
		v1.GET(`/ping1`, func(c contract.Context) {
			c.RenderWith(200, render.JSON, map[string]string{
				"message": "pong1",
			})
			c.Next()
		})
		v1.GET(`/for-run`, func(c contract.Context) {
			ctx2, _ := context.WithTimeout(c, time.Second*15)
			go func(ctx context.Context) {
			Next:
				for {
					select {
					case <-ctx.Done():
						fmt.Println(ctx.Err())
						break Next
					default:
						Add("firmeve")
					}
				}
				fmt.Println("execute over")
			}(ctx2)
			c.RenderWith(200, render.Text, "success")
		})
		v1.GET(`/panic`, func(c contract.Context) {
			panic(kernel.Error("something"))
			c.Next()
		})
	}
	debug := router.Group("/debug/pprof")
	{
		debug.Handler("GET", "/", pprof.Index)
		debug.Handler("GET", "/cmdline", pprof.Cmdline)
		debug.Handler("GET", "/profile", pprof.Profile)
		debug.Handler("GET", "/symbol", pprof.Symbol)
		debug.Handler("GET", "/trace", pprof.Trace)
	}
}

var datas = make([]string, 5500)

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	datas = append(datas, sData)

	return sData
}

func main() {
	firmeve.RunDefault(firmeve.WithConfigPath("../config.yaml"), firmeve.WithProviders(
		[]contract.Provider{
			new(App),
		},
	))
}
