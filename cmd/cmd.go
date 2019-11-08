package cmd

import (
	"github.com/firmeve/firmeve"
	"github.com/spf13/cobra"
)

type Commander interface {
	Cmd() *cobra.Command
}

func Root() *cobra.Command {
	//var version bool
	cmd := &cobra.Command{
		Use:     "firmeve",
		Short:   "Firmeve Framework",
		Version: firmeve.Version,
	}
	cmd.SetVersionTemplate(`{{with .Short}}{{printf "%s " .}}{{end}}{{printf "Version %s" .Version}}
`)

	return cmd
}
