## 简介

借鉴于`Laravel`思想，使用`Provider`来解决模块的依赖和解耦。

可以看到`Firmeve`中的模块都是使用`Provider`来注册和启动的。

在后续的开发的，我们推荐使用`Provider`来解决依赖和参数解耦。



## 创建服务提供者

### 注册方法

`Register`方法作为提供者注册初始化模块实例，所有的提供者在系统加载时都会优先使用注册方法，来完成需要的实例创建。

`Register`方法就是为了模块能够正常运行而提供的基础实例

### 启动方法

当所有的`Provider`加载并注册完成后，会自动调用`Boot`方法来唤醒当前模块。`Boot`方法就是模块启动的入口

### 基础示例

所有`Provider`都可以使用`kernel.BaseProvider`嵌入，`kernel.BaseProvider`主要是基于`Container`进行包装的快捷调用。

```go
package app

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `app`
}

// Register
func (p *Provider) Register() {
	config := p.Resolve(`config`).(*config2.Config).Item(`logging`)
	p.Bind(`logger`, New(config), container.WithShare(true))
}

// Boot method
func (p *Provider) Boot() {

}

```



### 挂载提供者

通常我们会在`main`函数中初始化时挂载我们需要的提供者

```go
func main()  {
	firmeve.Run(firmeve.WithProviders(
		[]contract.Provider{
			// base
			new(http.Provider),
			new(database.Provider),
			new(validator.Provider),
			new(jwt.Provider),

			// app
			new(providers.AppProvider),

			// component
			new(providers.SpiderProvider),
		},
	),firmeve.WithCommands(
		[]contract.Command{
			new(http.HttpCommand),
			new(spider.Command),
			new(database.MigrateCommand),
		},
	))
}
```

如果某个`Provider`需要依赖于另一个`Provider`才可以启动，那么我们可在服务内部进行注册，如：

```go
type AppProvider struct{
	kernel.BaseProvider
}

func (AppProvider) Name() string {
	return `app`
}

func (p *AppProvider) Register() {
  // 是否进行注册HttpProvider
  // 第二个参数表示，是否强制覆盖注册
  p.Firmeve.Register(new(HttpProvider), false)
}

func (p *AppProvider) Boot() {
	
}

```

