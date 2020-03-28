package contract

import "github.com/spf13/cobra"

type (
	BaseCommand interface {
		Application() Application

		Providers() []Provider

		AddCommand(commands ...Command)

		Run() error

		Root() *cobra.Command

		Resolve(abstract interface{}, params ...interface{}) interface{}
	}

	Command interface {
		CobraCmd() *cobra.Command

		Run(root BaseCommand, cmd *cobra.Command, args []string)
	}
)
