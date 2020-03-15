package contract

import "github.com/spf13/cobra"

type Command interface {
	Cmd() *cobra.Command

	SetApplication(app Application)

	SetProviders(providers []Provider)

	Application() Application

	Providers() []Provider
}
