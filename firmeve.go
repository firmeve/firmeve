package firmeve

import (
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	logging "github.com/firmeve/firmeve/logger"
	"github.com/firmeve/firmeve/support"
)

var (
	initProviders = []contract.Provider{
		new(event.Provider),
		new(logging.Provider),
	}

	defaultProviders = []contract.Provider{
		new(http.Provider),
	}

	defaultCommands = []contract.Command{
		new(http.HttpCommand),
	}
)

type (
	option struct {
		providers []contract.Provider
		commands  []contract.Command
	}
)

func WithProviders(providers []contract.Provider) support.Option {
	return func(object support.Object) {
		object.(*option).providers = providers
	}
}

func WithCommands(commands []contract.Command) support.Option {
	return func(object support.Object) {
		object.(*option).commands = commands
	}
}

func RunDefault(options ...support.Option) (contract.Application, contract.BaseCommand) {
	option := parseOption(options)
	option.providers = append(defaultProviders, option.providers...)
	option.commands = append(defaultCommands, option.commands...)

	return Run(WithProviders(option.providers), WithCommands(option.commands))
}

func Run(options ...support.Option) (contract.Application, contract.BaseCommand) {
	option := parseOption(options)

	// init providers
	option.providers = append(initProviders, option.providers...)

	command := kernel.NewCommand(option.providers, option.commands...)

	command.Run()

	return command.Application(), command
}

func parseOption(options []support.Option) *option {
	return support.ApplyOption(&option{
		providers: make([]contract.Provider, 0),
	}, options...).(*option)
}
