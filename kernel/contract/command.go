package contract

import "github.com/spf13/cobra"

type (
	BaseCommand interface {
		//SetApplication(app Application)

		//SetProviders(providers []Provider)

		Application() Application

		Providers() []Provider

		AddCommand(cmds ...Command)

		Run() error
	}

	Command interface {
		Cmd() *cobra.Command

		Run(root BaseCommand, cmd *cobra.Command, args []string)
		//Name() string
		//Run()

	}
)
