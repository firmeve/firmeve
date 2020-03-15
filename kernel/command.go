package kernel

import (
	"github.com/fatih/color"
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
	Command  *cobra.Command
}

func (c *Command) SetProviders(providers []contract.Provider) {
	c.Provider = providers
}

func (c *Command) Providers() []contract.Provider {
	return c.Provider
}

func (c *Command) SetApplication(app contract.Application) {
	c.Firmeve = app
}

func (c *Command) Application() contract.Application {
	return c.Firmeve
}

func CommandRoot(app contract.Application) *cobra.Command {
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
