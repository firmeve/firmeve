## 简介
Firmeve command 基于 cobra 命令封装，可以轻松集成。

## 编写命令

增加Firmeve命令只需要实现`kernel.Command` 接口

```go
Command interface {
    CobraCmd() *cobra.Command

    Run(root BaseCommand, cmd *cobra.Command, args []string)
}
```

在系统入口入增加命令
```go
firmeve.Run(
    firmeve.WithCommands([]contract.Command{
        new(Command)
    }),
)
```

## 示例命令
```go
type ExampleCommand struct {
	command *cobra.Command
}

func (c *ExampleCommand) CobraCmd() *cobra.Command {
	command := new(cobra.Command)
	command.Use = "example"
	command.Short = "A example demo"

	return command
}

func (c *ExampleCommand) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
    //do something
}

```

> 更多关于`Flags`和`Params`参考 [https://github.com/spf13/cobra](https://github.com/spf13/cobra)