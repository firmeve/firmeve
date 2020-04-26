## 简介

配置模块是Firmeve中基础模块之一，Firmeve的启动必须带上`config`配置，在cli中使用`-c`选项指定`config`路径



## 配置文件示例

```yaml
framework:
  lang: zh-CN
  key: "!!@#$123^%"
database:
  default: mysql
  connections:
    mysql:
      addr: "root:root@(127.0.0.1)/default?charset=utf8mb4&parseTime=True&loc=Local"
    pool:
      max_idle: 100
      max_connection: 50
      max_lifetime: 60
  migration:
    path: "../../../testdata/migrations"

cache:
  prefix: firmeve_cache
  default: redis
  repositories:
    redis:
      connection: cache
```



## 基础示例

```go
// 创建config
var c = config.New()

// 调用节点，如上配置示例
databaseConfig := c.Item("database")

// 基本用法
databaseConfig.Get("default")

// 多层级调用
databaseConfig.Get("connections.mysql.addr")

// 判断一个key是否存在
if databaseConfig.Has("default") {
  fmt.Println("exists")
}
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



## 环境变量

环境变量的键会统一转换为大写，**不支持小写**环境变量名

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