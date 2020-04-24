## 简介
日志模块是`Firmeve`框架最基础的模块之一，在默认系统启动时就会自动加载。

`Firmeve`日志是继承于`uber zap`。

`Firmeve`日志目前支持`console`和`file`以及`stack`混合等3种类型。



## 配置

```yaml
# 默认日志信道
default: stack
channels:
  stack:
    - file
    - console
  # 文件日志
	file:
		# 日志文件路径
    path: "../../../testdata/logs"
    # megabytes
    size:    100
    # 最大备份天数
    backup: 3
    age:     1
    # 日志级别
    level: debug
  
  # 控制台日志 os.Stdout
	console:
		level: debug

# 记录消息栈级别
stack_level: error
```



## 基础示例

### 创建日志实例

```go
loggerConfig := config.(*config2.Config).Item(`logging`)
var logger = logger.New(loggerConfig)
```

> config 请参见config模块配置获取



### 基础用法

```go
// debug
logger.Debug("Debug info")

// 附加参数
// {"level":"DEBUG","time":"2020-04-23 13:17:59","message":"Debug info","error":"something","stacktrace":"...."}
logger.Debug("Debug info" , "error" , errors.New("something")) 

logger.Warn("Warn info")

```



### 日志级别

- debug
- info
- warn
- error
- fatal
- panic



### 支持的类型

| 类型    | 说明                    |
| :------ | :---------------------- |
| console | 直接输出到console的日志 |
| file    | 写入到指定的文件日志    |
| stack   | 多种类型通道合并记录    |
