## 简介
`Firmeve` Config



## 基础示例

### 创建一个`Config`
```go
directory := "./testdata/config/config.yaml"
config := config.New(directory)
```

### Config Item
在`Config`内部每个扫描的文件名（去除后缀名）作为`item`的调用名称，如下
`Item`是`Configurator Interface`的实现
```go
config.Item("app")
```
> 注意：目前只支持yaml格式配置文件

### 值的设置
```go
// 设置一个默认值
config.Item("app").SetDefault("key", "value")
// 修改或添加一个值
config.Item("app").Set("key", "value")
```

> 虽然提供了修改值的接口，但强烈建议不使用该方法
> http.Context内不可直接使用Set来修改值，多个goroutine并发会导致不可知问题
> 如果确实需要修改，请copy一份config副本，并后续使用此副本

### 值的获取
```go
config.Item("app").Get("key")
// 多层级调用
config.Item("app").Get("key.foo.bar")
```
其它可用方法：
- GetBool()
- GetString()
- GetFloat64()
- GetInt()
- GetIntSlice()
- GetString()
- GetStringMap()
- GetStringMapString()
- GetTime()
- GetDuration()

### 判断Key是否存在 
```go
config.Item("app").Exists("key")
```

### 加载一个Item
```go
file := "./config/tmp.yaml"
config.Load(file)
```

## 环境变量

使用`Config`环境变量的键会统一转换为大写，**不支持小写**环境变量名

### 环境变量设置
```go
// go自动环境变量的设置
os.Setenv("FOO", "foo")
// config 当使用SetEnv函数设置
SetEnv("bar", "bar")
// 
os.Setenv("baz","baz")
``` 
### 环境变量读取
```go
// 可以读取
fmt.Println(GetEnv("FOO"))
// 可以读取
fmt.Println(GetEnv("bar"))
fmt.Println(os.Getenv("BAR"))
// 无法读取
fmt.Println(os.Getenv("baz"))
```