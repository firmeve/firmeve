## 简介
`Database`模块是基于`gorm`的统一扩展和封装

## 基础示例

### 新建连接
```go
config := Firmeve.Get(`config`).(*config.Config).Item(`database`)
DB := New(config)
```

### 默认连接
```go
// 新的连接
DB.Connection(driver string)
// 默认连接
DB.ConnectionDefault()
```

### 关闭连接

```go
DB.Close(driver)
DB.CloseDefault()
```