package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Commander interface {
	Cmd() *cobra.Command
}

func New(version string) *cobra.Command {
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
	cmd := &cobra.Command{
		Use:     `firmeve`,
		Short:   `Firmeve Framework [` + version + `]`,
		Long:    logoColor.Sprint(logo) + `Framework [` + versionColor.Sprint(version) + `]`,
		Version: version,
	}
	//cmd.PersistentFlags().StringP("config", "C", "", "Config directory path(required)")
	//err := cmd.MarkFlagRequired("config")
	//if err != nil {
	//	firmeve.F(`logger`).(logging.Loggable).Fatal(err.Error())
	//}
	cmd.SetVersionTemplate("{{with .Short}}{{printf \"%s \" .}}{{end}}{{printf \"Version %s\" .Version}}\n")

	return cmd
}
