package kernel

import (
	"github.com/fatih/color"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
	"github.com/spf13/cobra"
	"os"
)

const (
	defaultConfigPath2 = `../testdata/config`
)

type Command struct {
	Firmeve  contract.Application
	Provider []contract.Provider
	//commands []*cobra.Command
	root *cobra.Command
}

//@todo 暂时这样，后面可以做移除命令
//@todo 而且这个还有问题，命令可以重复添加，可以考虑reflect.Type惟一性
func (c *Command) AddCommand(cmds ...contract.Command) {
	for _, v := range cmds {
		cmd := v.Cmd()
		cmd.Run = func(cmd *cobra.Command, args []string) {
			// bootstrap
			// kernel.BootFromCommand(c, cmd)
			configPath := cmd.Flag(`config`).Value.String()
			devMode := cmd.Flag(`dev`).Value.String()
			devModeBool := false
			if devMode == `true` {
				devModeBool = true
			}

			//kernel.Boot(configPath, devModeBool, root.Application(), root.Providers())
			c.Boot(configPath, devModeBool)

			v.Run(c, cmd, args)
		}

		c.root.AddCommand(cmd)
		//c.commands = append(c.commands, cmd)
	}
}
func (c *Command) Boot(configPath string, devMode bool) contract.Application {
	var mode uint8
	if devMode {
		mode = contract.ModeDevelopment
	} else {
		mode = contract.ModeProduction
	}
	c.Firmeve.SetMode(mode)

	c.Firmeve.Bind("firmeve", c.Firmeve)

	c.Firmeve.Bind(`config`, config.New(configPath), container.WithShare(true))

	//registerBaseProvider(app)
	//if providers != nil && len(providers) != 0 {
	c.Firmeve.RegisterMultiple(c.Provider, false)
	//}

	c.Firmeve.Boot()

	return c.Firmeve
}

//func (c *Command) SetProviders(providers []contract.Provider) {
//	c.Provider = providers
//}

func (c *Command) Run() error {
	return c.root.Execute()
}

func (c *Command) Providers() []contract.Provider {
	return c.Provider
}

//func (c *Command) SetApplication(app contract.Application) {
//	c.Firmeve = app
//}

func (c *Command) Application() contract.Application {
	return c.Firmeve
}

func NewCommand(providers []contract.Provider, commands ...contract.Command) contract.BaseCommand {
	app := New()
	command := &Command{
		Firmeve:  app,
		Provider: providers,
		//commands: make([]*cobra.Command, 0),
		root: rootCommand(app),
	}
	command.AddCommand(commands...)
	return command
}

func rootCommand(app contract.Application) *cobra.Command {
	logo := `
 _____   _   _____        ___  ___   _____   _     _   _____
|  ___| | | |  _  \      /   |/   | | ____| | |   / / | ____|
| |__   | | | |_| |     / /|   /| | | |__   | |  / /  | |__
|  __|  | | |  _  /    / / |__/ | | |  __|  | | / /   |  __|
| |     | | | | \ \   / /       | | | |___  | |/ /    | |___
|_|     |_| |_|  \_\ /_/        |_| |_____| |___/     |_____|

`
	logoColor := color.New(color.FgCyan, color.Bold)
	versionColor := color.New(color.FgRed, color.Bold)
	version := app.Version()
	cmd := &cobra.Command{
		Use:     `firmeve`,
		Short:   `Firmeve Framework [` + version + `]`,
		Long:    logoColor.Sprint(logo) + `Framework [` + versionColor.Sprint(version) + `]`,
		Version: version,
	}

	cmd.PersistentFlags().StringP("config", "c", getConfigPath(path.RunRelative(defaultConfigPath2)), "config path")
	cmd.PersistentFlags().BoolP("dev", "", false, "open development mode (default production)")

	cmd.SetVersionTemplate("{{with .Short}}{{printf \"%s \" .}}{{end}}{{printf \"Version %s\" .Version}}\n")

	cmd.SetArgs(os.Args[1:])
	return cmd
}

func getConfigPath(path string) string {
	if path == `` {
		return os.Getenv("FIRMEVE_CONFIG")
	}

	return path
}
