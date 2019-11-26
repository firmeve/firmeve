## 简介
应用系统中缓存必不可少
，`Firmeve`提供了一套，易于外部调用以及扩展

## 基础示例

### 创建一个新的缓存资源
```go
// 读取config中cache文件配置
var config config.Configurator
cacheConfig = config.New(path string).Item(`cache`)
cache := New(cacheConfig)
```

> 具体请参考[cacheProvider](https://github.com/firmeve/firmeve/blob/develop/cache/provider.go)的启动方式

### 获取当前缓存实例
```go
Firmeve.Get(`config`).(*config.Config)
```

### 增加数据缓存
```go
// 普通数据
if err := cache.Put(`key`, `value`, time.Now().Add(time.Second*1000)); err != nil {
	//错误处理
}

// 需要序列化的数据
v := struct {
    Id   int
    Name string
}{
    Id:   10,
    Name: "simon"
}

if err := cache.PutEncode(`key`, v, time.Now().Add(time.Hour)); err != nil {
	//错误处理
}

// 永久存储
if err := cache.Forever(`key`, `value`); err != nil {
	//错误处理
}

if err := cache.ForeverEncode(`key`, `value`); err != nil {
	//错误处理
}
```

### 获取数据缓存
```go
// 普通数据
if err := cache.Get(`key`); err != nil {
	//错误处理
}

// 获取值，如果不存在则返回默认值
if err := cache.GetDefault(`key`, `value`); err != nil {
	//错误处理
}

// 取出值，如果存在则删除
if err := cache.Pull(`key`); err != nil {
	//错误处理
}

// 取出值，如果不存在则返回默认值，如果存在则删除
if err := cache.PullDefault(`key`); err != nil {
	//错误处理
}
```

### 删除缓存
```go
// 删除指定key
if err := cache.Forget(`key`); err != nil {
	//错误处理
}

// 清除所有缓存
if err := cache.Flush(); err != nil {
	//错误处理
}
```

### 其它接口
```go
Get(key string) (interface{}, error)

Add(key string, value interface{}, expire time.Time) error

Put(key string, value interface{}, expire time.Time) error

Forever(key string, value interface{}) error

Forget(key string) error

Increment(key string, steps ...int64) error

Decrement(key string, steps ...int64) error

Has(key string) bool

Flush() error

GetDefault(key string, defaultValue interface{}) (interface{}, error)

Pull(key string) (interface{}, error)

PullDefault(key string, defaultValue interface{}) (interface{}, error)

GetDecode(key string, to interface{}) (interface{}, error)

AddEncode(key string, value interface{}, expire time.Time) error

ForeverEncode(key string, value interface{}) error

PutEncode(key string, value interface{}, expire time.Time) error
```

### 扩展驱动

目前只支持`Redis`驱动，如果需要扩展驱动也十分方便

```go
cache.Register(driverName, store repository.Cacheable)
```

需要实现的接口
```go
Cacheable interface {
    Get(key string) (interface{}, error)

    Add(key string, value interface{}, expire time.Time) error

    Put(key string, value interface{}, expire time.Time) error

    Forever(key string, value interface{}) error

    Forget(key string) error

    Increment(key string, steps ...int64) error

    Decrement(key string, steps ...int64) error

    Has(key string) bool

    Flush() error
}
```