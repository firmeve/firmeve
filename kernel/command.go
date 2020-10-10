package kernel

import (
	"github.com/fatih/color"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/cobra"
	"os"
)

type (
	command struct {
		application contract.Application
		providers   []contract.Provider
		root        *cobra.Command
		mount       func(app contract.Application)
	}

	CommandOption struct {
		ConfigPath string
		Providers  []contract.Provider
		Commands   []contract.Command
		Mount      func(app contract.Application)
	}
)

func (c *command) AddCommand(commands ...contract.Command) {
	for i := range commands {
		v := commands[i]
		cmd := v.CobraCmd()
		cmd.Run = func(cmd *cobra.Command, args []string) {
			// init params --> *important*
			configPath := cmd.Flag(`config`).Value.String()
			devMode := cmd.Flag(`dev`).Value.String()
			devModeBool := false
			if devMode == `true` {
				devModeBool = true
			}

			// bootstrap
			c.boot(configPath, devModeBool)
			c.debugLog(cmd.Use + ` command register and booted`)

			// mount
			if c.mount != nil {
				c.mount(c.Application())
				c.debugLog(`Application func mounted`)
			}

			// panic handler
			defer Recover(c.Resolve(`logger`).(contract.Loggable))

			v.Run(c, cmd, args)
		}

		c.root.AddCommand(cmd)
	}
}

func (c *command) boot(configPath string, devMode bool) {
	var mode uint8
	if devMode {
		mode = contract.ModeDevelopment
	} else {
		mode = contract.ModeProduction
	}
	c.application.SetMode(mode)

	c.application.Bind(`application`, c.application)
	c.application.Bind(`config`, NewConfig(configPath), container.WithShare(true))

	c.application.RegisterMultiple(c.providers, false)

	c.application.Boot()
}

func (c *command) debugLog(context ...interface{}) {
	c.application.Resolve(`logger`).(contract.Loggable).Debug(context...)
}

func (c *command) Run() error {
	return c.root.Execute()
}

func (c *command) Root() *cobra.Command {
	return c.root
}

func (c *command) Providers() []contract.Provider {
	return c.providers
}

func (c *command) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return c.Application().Make(abstract, params...)
}

func (c *command) Application() contract.Application {
	return c.application
}

func NewCommand(option *CommandOption) contract.BaseCommand {
	app := New()
	command := &command{
		application: app,
		providers:   option.Providers,
		root:        rootCommand(app, option.ConfigPath),
		mount:       option.Mount,
	}
	command.AddCommand(option.Commands...)
	return command
}

func rootCommand(app contract.Application, configPath string) *cobra.Command {
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

	cmd.PersistentFlags().StringP("config", "c", getConfigPath(configPath), "config path")
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
