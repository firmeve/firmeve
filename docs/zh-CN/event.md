## 简介

事件模块是`Firmeve`的基础核心模块之一，事件模式可以让你的代码轻松解偶，更主要的是事件钩子可以让你轻松控制运行流程。



## 基本用法

任何事件都必须实现`contract.EventHandler`接口

```go
type EventHandler interface {
  Handle(params ...interface{}) (interface{}, error)
}
```



```go
// 
var e = event.New()

// 注册事件
e.Listen("foo", EventHandler)

// 注册多个事件
e.ListenMany("foo", []EventHandler)

// 调度事件
e.Dispatch("foo", "params1", "params2")
```

> 注意：
>
> 当事件中的函数执行失败后(err != nil)，后续函数将不会执行，但会返回之前执行的全部结果