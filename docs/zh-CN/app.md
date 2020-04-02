## 简介
`Firmeve` 是一个提供了多基础组件的web框架

## 基础架构

![base](../images/base.png)

## 生命周期

### Main
在 Main 中我们需要执行`firmeve.Run`或者`firmeve.RunDefault`函数，并且进行 Provider 和 Command 挂载

### RootCommand
一切都是从RootCommand创建后开始
RootCommand中会首先启动一个新的应用容器kernel.New()和一个新的Root Command
并且会注册main中加入的provider和sub command，注册命令时，会自动包装所有子命令的Run方法

### Run包装
获取RootCommand的运行配置文件，以及运行模式等基础参数，并进行boot启动加载

### Boot
在Boot方法中，会设置应用运行模式，绑定基础实例，挂载Provider以及启动已存在的Provider

### SubCommand
最后在各子命令中执行需要运行的最后方法


## 基础示例

在`main.go`中增加此代码

```go
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
	firmeve.RunDefault(firmeve.WithProviders(
		[]contract.Provider{
			new(App),
		},
	))
}
```
运行启动命令
```bash
go run main.go http:serve
```