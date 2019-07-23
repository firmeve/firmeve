## 简介
借鉴于`Laravel`思想，使用`Provider`来解决模块的依赖和解耦

## 基础示例

```go
import (
	"github.com/firmeve/firmeve"
)

type CacheServiceProvider struct {
	firmeve.FirmeveServiceProvider
}

func (csp *CacheServiceProvider) Boot() {

}

func (csp *CacheServiceProvider) Register() {
	// 绑定CacheManager
	csp.Firmeve.Bind(`cache`, NewManager, firmeve.WithBindShare(true))
}
```

## Provider挂载
```go
firmeve.GetFirmeve().Register(`cache`, new(CacheServiceProvider))
```

## Register方法

在`Register`方法中我们通常使用容器来实现实例绑定，如：

```go
csp.Firmeve.Bind(`cache`, NewManager, firmeve.WithBindShare(true))
```

## Boot方法
**Boot方法在所有服务提供者被注册以后才会被调用**
