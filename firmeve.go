package firmeve

import (
	config2 "github.com/firmeve/firmeve/config"
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

	Application contract.Application
	Logger      contract.Loggable
	Config      *config2.Config
	Event       contract.Event
)

type (
	option struct {
		providers  []contract.Provider
		commands   []contract.Command
		configPath string
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

func WithConfigPath(path string) support.Option {
	return func(object support.Object) {
		object.(*option).configPath = path
	}
}

func RunDefault(options ...support.Option) error {
	option := parseOption(options)
	option.providers = append(defaultProviders, option.providers...)
	option.commands = append(defaultCommands, option.commands...)

	return Run(WithConfigPath(option.configPath), WithProviders(option.providers), WithCommands(option.commands))
}

func Run(options ...support.Option) error {
	return RunWithFunc(nil, options...)
}

func RunWithFunc(f func(application contract.Application), options ...support.Option) error {
	option := parseOption(options)

	// init providers
	option.providers = append(initProviders, option.providers...)

	command := kernel.NewCommand(&kernel.CommandOption{
		ConfigPath: option.configPath,
		Providers:  option.providers,
		Commands:   option.commands,
		Mount:      f,
	})

	return command.Run()
}

func RunWithSupportFunc(f func(application contract.Application), options ...support.Option) error {
	return RunWithFunc(func(application contract.Application) {
		bindingBaseService(application)
		if f != nil {
			f(application)
		}
	}, options...)
}

func bindingBaseService(app contract.Application) {
	Application = app
	Logger = app.Resolve(`logger`).(contract.Loggable)
	Config = app.Resolve(`config`).(*config2.Config)
	Event = app.Resolve(`event`).(contract.Event)
}

func parseOption(options []support.Option) *option {
	return support.ApplyOption(&option{
		providers: make([]contract.Provider, 0),
	}, options...).(*option)
}
