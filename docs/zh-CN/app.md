## 简介
`Firmeve` 是一个提供了多基础组件的web框架

## 基础示例

在`main.go`中增加此代码

```go
package main

import (
	"os"
	path2 "github.com/firmeve/firmeve/support/path"
	"github.com/firmeve/firmeve/bootstrap"
	cmd2 "github.com/firmeve/firmeve/http/cmd"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cmd"
	"github.com/firmeve/firmeve/http"
)

func main() {
	// Base bootstrap
	app := firmeve.New()
	bootstrap2 := bootstrap.New(app, bootstrap.WithConfigPath(path2.RunRelative(`../../testdata/config`)))
	bootstrap2.FastBootFull()

    router := app.Get(`http.router`).(*http.Router)
	v1 := router.Group("/api/v1")
	{
		v1.GET(`/ping`, func(c *http.Context) {
			c.Data(map[string]string{
				"message": "pong",
			})
			c.Next()
		})
	}
	// Command
	root := cmd.Root()
	root.AddCommand(cmd2.NewServer(bootstrap2).Cmd())
	root.SetArgs(os.Args[1:])
	root.Execute()
}
```
运行启动命令
```bash
go run main.go http:serve
```