## 简介
`Firmeve`提供一了个简洁的`Event`事件绑定

## 基本用法

### 创建事件
```go
event := event.New()
```

### 事件接口

事件的处理需要实现`Handler`接口

```go
type Handler interface {
    Handle(params ...interface{}) (interface{}, error)
}

```

### 注册事件
```go
// 注册单个事件
event.Listen("foo", Handler)

// 注册多个事件
event.ListenMany("foo", []Handler)

```

### 触发事件
```go
if event.Has("foo") {
    results := event.Dispatch("foo")
    fmt.Println(results)	
}
```

> 注意：当事件中的函数执行失败后(err != nil)，后续函数将不会执行，但会返回之前执行的全部结果

## 自定义事件接口
```go
IDispatcher interface {
    Listen(name string, df Handler)
    ListenMany(name string, df handlers)
    Dispatch(name string, params ...interface{}) []interface{}
    Has(name string) bool
}
```