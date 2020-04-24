```
## 简介
`Database`模块是基于`gorm`的统一扩展和封装

## 基础示例

### 新建连接
​```go
config := Firmeve.Get(`config`).(*config.Config).Item(`database`)
DB := New(config)
​```

### 默认连接
​```go
// 新的连接
DB.Connection(driver string)
// 默认连接
DB.ConnectionDefault()
​```

### 关闭连接

​```go
DB.Close(driver)
DB.CloseDefault()
​```
```

## 简介

`Firmeve Database` 模块是基于`gorm`的统一扩展和封装。



## 基础配置

```yaml
database:
	# 默认连接
  default: mysql
  # 数据库连接
  connections:
    mysql:
      addr: "root:root@(127.0.0.1)/default?charset=utf8mb4&parseTime=True&loc=Local"
    # 连接池设置  
    pool: 
      max_idle: 100
      max_connection: 50
      max_lifetime: 60
  # 数据库迁移配置
  migration:
    path: "../../../testdata/migrations"
```

> 其它如pgsql,sqlite，请参考gorm



### 注册Database模块

```go
func main() {
	firmeve.RunDefault(firmeve.WithProviders(
		[]contract.Provider{
			new(database.Provider),
		},
	), firmeve.WithCommands(
		[]contract.Command{
			new(database.MigrateCommand),
		},
	))
}
```





## 基础示例

```go
// 获取当前默认连接
var conn1 = p.Resolve(`db.connection`).(*gorm.DB)

// 获取指定连接
var conn2 = p.Resolve(`db`).(*database.DB).Connection("sqlite").(*gorm.DB)

//关闭连接
conn.close()
conn2.close()
```



## 数据表迁移

如上示例所示，执行`migration`命令需要提前挂载到`root`根命令下，`firmeve migration`主要是基于`golang-migrate/migrate`库开发的命令扩展。

默认只支持`mysql`数据库自动迁移，如果需要支持pgsql或其它数据库请自动引入相关数据库

```go
import (
  _ "github.com/golang-migrate/migrate/database/postgres"
  _ "github.com/golang-migrate/migrate/database/sqlite3"
)
```

主要命令如下：

```bash
## 查看所有命令
go run main.go -c .config.yaml migrate --help

# 创建迁移
migrate create NAME

# 执行迁移步骤
migrate step N 

# 回滚步骤
migrate rollback N    

# 运行所有迁移
migrate up     

# 撤销所有迁移
migrate down      

# 强制忽略失败的迁移
migrate force V

# 当前迁移版本
migrate version

```

更多示例请参考：[database/command](https://github.com/firmeve/firmeve/examples/database/migrate)