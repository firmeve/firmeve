# This is a web framework that firmly believes in dreams.

```
Firmeve = Firm + Believe
``` 

Be firm in your will and believe in your ideals.(坚定自己的意志，相信自己的理想。)

Those who have achieved nothing can always tell you that you can't make a big deal. If you have an ideal, you have to defend it.(那些一事无成的人总是告诉你，你也成不了大器，如果你有理想的话，就要去捍卫它。)


[![Build Status](https://travis-ci.com/firmeve/firmeve.svg?branch=develop)](https://travis-ci.com/firmeve/firmeve)
[![codecov](https://codecov.io/gh/firmeve/firmeve/branch/develop/graph/badge.svg)](https://codecov.io/gh/firmeve/firmeve)
[![GitHub license](https://img.shields.io/github/license/firmeve/firmeve.svg)](https://github.com/firmeve/firmeve/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/firmeve/firmeve)](https://goreportcard.com/report/github.com/firmeve/firmeve)

## Quick start
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

## Feature list
- **Core**
    - [x] [Ioc](./docs/zh-CN/container.md)
    - [x] [Application](./docs/zh-CN/app.md)
    - [x] [Config](./docs/zh-CN/config.md)
    - [x] [Provider](./docs/zh-CN/provider.md)
    - [x] [Event](./docs/zh-CN/event.md)
    - [x] [Logger](./docs/zh-CN/logger.md)
    - [ ] [Command](./docs/zh-CN/command.md)
- **Base**
    - [x] [Database](./docs/zh-CN/databases.md)
    - [x] [Cache](./docs/zh-CN/cache.md)
    - [ ] [Queue](./docs/zh-CN/queue.md)
    - [ ] [Cron](./docs/zh-CN/cron.md)
- **Extension**
