## 简介
系统日志，所需要配置参见：`testdata/config/logging.yaml`

## 基础示例

### 获取实例
```go
// Get firmeve container logger
logger := firmeve.Instance().Get(`logger`)
// Create a new logger
logger := logger.New(&Config{
    Current: `console`,
    Channels: ConfigChannelType{
        `stack`: []string{`file`, `console`},
        `console`: ConfigChannelType{
            `level`: `debug`,
        },
        `file`: ConfigChannelType{
            `level`:  `debug`,
            `path`:   "../testdata/logs",
            `size`:   100,
            `backup`: 3,
            `age`:    1,
        },
    },
})
// Or
logger := logger.Default()
```

### 基础用法
```go
logger.Debug("message")

logger.Error("message")

// 传入附加的参数
logger.Debug("message", "url", "https://firmeve.com")
```

### 指定日志通道
```go
channel := logger.Channel(`console`)
```
> 注意：此时的新日志通道并不是原日志记录通道

### 支持的类型
| 类型 | 说明 |
| :-----| :---- |
| console | 直接输出到console的日志 |
| file | 写入到指定的文件日志 |
| stack | 多种类型通道合并记录 |

### 错误级别
- Debug
- Info
- Notice
- Warn
- Error
- Fatal
- ~~Critical~~
- ~~Alert~~
- ~~Emergency~~