package firmeve

import (
	"github.com/firmeve/firmeve/cmd"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/support"
	"github.com/spf13/cobra"
	"os"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/logger"
)

type (
	Runner interface {
		Run()
	}

	Firmeve struct {
		kernel.IApplication
	}

	option struct {
		providers []kernel.IProvider
	}
)

func WithProviders(providers []kernel.IProvider) support.Option {
	return func(object support.Object) {
		object.(*option).providers = providers
	}
}

func New(mode uint8, configPath string, options ...support.Option) Runner {
	return boot(mode,configPath,parseOption(options)).(Runner)
}

func Default(mode uint8, configPath string, options ...support.Option) Runner {
	defaultProviders := []kernel.IProvider{
		new(http.Provider),
	}

	option := parseOption(options)
	option.providers = append(defaultProviders,option.providers...)

	return boot(mode, configPath,option).(Runner)
}


func boot(mode uint8, configPath string, option *option) kernel.IApplication {
	f := &Firmeve{
		kernel.New(mode),
	}

	f.Bind("firmeve", f)

	f.Bind(`config`, config.New(configPath), container.WithShare(true))

	f.registerBaseProvider()

	if len(option.providers) != 0 {
		f.RegisterMultiple(option.providers, false)
	}

	f.Boot()

	return f
}

func parseOption(options []support.Option) *option {
	return support.ApplyOption(&option{
		providers: make([]kernel.IProvider, 0),
	}, options...).(*option)
}

func (f *Firmeve) Run() {
	root := f.Get("command").(*cobra.Command)
	root.SetArgs(os.Args[1:])
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}

func (f *Firmeve) registerBaseProvider() {
	f.RegisterMultiple([]kernel.IProvider{
		new(cmd.Provider),
		new(event.Provider),
		new(logging.Provider),
	}, false)
}