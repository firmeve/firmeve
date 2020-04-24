## 简介

`Firmeve`命令基于`cobra`，在其基础之上提供了一系内置可用命令，以及灵活的命令扩展。命令是`Firmeve`运行的最初入口，`Firmeve`是通过多命令运行。



## 基础示例



### 创建命令

命令的创建必须实现`contract.Command`接口

```go
Command interface {
  // 返回一个新的 cobra command
  CobraCmd() *cobra.Command

  // 业务主入口
  Run(root BaseCommand, cmd *cobra.Command, args []string)
}
```

以下是一个基础示例：

```go
import (
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/cobra"
)

type TestCommand struct {
}

func (t TestCommand) CobraCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "testing"
	cmd.Short = "Testing a cmd"
	return cmd
}

func (t TestCommand) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	fmt.Println("run")
}
```

> 更多用法，请参见示例 [example/command](https://github.com/firmeve/firmeve/examples/command/main.go)

### 命令注册

所有的命令都是在`main`初始化中注册，本质上来说，是所有的子命令会在`main`中挂载到系统`root command`中

```go
func main() {
	firmeve.RunDefault(firmeve.WithCommands(
		[]contract.Command{
			new(TestCommand),
		},
	))
}
```



### 内置命令

Firmeve提供了一系列内置运行命令，如：`http:serve` `migrate`等，具体可执行`go run main.go -c config.yaml`查看



## 更多用法

关于更多命令开发方法请参考 [cobra](https://github.com/spf13/cobra)



