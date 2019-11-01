## 简介
`Firmeve`提供一了个简洁的`Event`事件绑定

## 基本用法

### 注册事件
```go
dispatcher := event.NewDispatcher()
// or 
f := New()
dispatcher := f.Get(`dispatcher`).(*event.Dispatcher)

dispatcher.listen("eventName", func (params ...interface{}) interface{}{
	...
})
```

### 触发事件
```go
results := dispatcher.dispatch("eventName")
fmt.Println(results)
```

> 注意：当事件中的函数执行失败后，后续函数将不会执行，但会返回之前执行的全部结果