## 简介

为什么要开发这个框架？

强如`gin`,`iris`这类web框架，虽然性能非常但是在开发应用过程中我们还是要使用大量的第三方包封装组合使用，所以才有的Firmeve的作用，`Firmeve`中包装整合了常用的第三方稳定包，并且提供了快速集成规范。同时本身也提供了稳定的框架底层基础，让开发者可以轻松快速上手和集成任何第三方扩展。



## 基础示例

```go
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

```

运行`main`

```bash
go run main.go -c config.yaml http:serve
```



## 基础架构

![base](../images/base.png)



## 生命周期

todo ....



## 后续开发计划

