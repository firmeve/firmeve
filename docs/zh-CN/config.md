## 简介
`Firmeve` Config

## Config Item
在`Config`内部每个扫描的文件名（去除后缀名）作为`item`的调用名称，如
```go
GetConfig().Item("app")
```
> 注意：目前只支持yaml格式配置文件

## 基础示例

### 创建一个`Config`
```go
directory := "./testdata/config"
config := NewConfig(directory)
```
当创建完成后可以调用`GetConfig()来调用`

> 注意：Config是一个单例对象


### 值的设置
```go
// 设置一个默认值
GetConfig().Item("app").SetDefault("key", "value")
// 修改或添加一个值
GetConfig().Item("app").Set("key", "value")
```

### 值的获取
```go
GetConfig().Item("app").Get("key")
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
GetConfig().Item("app").Exists("key")
```

### 加载一个Item
```go
file := "./config/tmp.yaml"
GetConfig().Load(file)
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