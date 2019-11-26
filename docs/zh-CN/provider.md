## 简介
借鉴于`Laravel`思想，使用`Provider`来解决模块的依赖和解耦

## 基础示例

```go
import (
	"github.com/firmeve/firmeve"
)

type Provider struct {
	firmeve.BaseProvider
}

func (p *Provider) Name() string {
	return `app`
}

func (p *Provider) Register() {
	// bind
	p.Firmeve.Bind(`foo`, func() string {
		return `bar`
	})
}

func (p *Provider) Boot() {
}
```

> 注意：Boot方法在所有服务提供者被注册以后才会被调用

## Provider挂载
```go
f := firmeve.New()
f.Register([]firmeve.Provider{
    new(Provider),
})
```
