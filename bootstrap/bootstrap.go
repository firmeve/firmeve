package kernel

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/kernel/contract"
	logging "github.com/firmeve/firmeve/logger"
)

func BootFromCommand(cmd contract.Command) contract.Application {
	configPath := cmd.Cmd().Flag(`config`).Value.String()
	devMode := cmd.Cmd().Flag(`dev`).Value.String()
	devModeBool := false
	if devMode == `true` {
		devModeBool = true
	}

	return Boot(configPath, devModeBool, cmd.Application(), cmd.Providers())
}

func Boot(configPath string, devMode bool, app contract.Application, providers []contract.Provider) contract.Application {
	var mode uint8
	if devMode {
		mode = contract.ModeDevelopment
	} else {
		mode = contract.ModeProduction
	}
	app.SetMode(mode)

	app.Bind("firmeve", app)

	app.Bind(`config`, config.New(configPath), container.WithShare(true))

	registerBaseProvider(app)
	if providers != nil && len(providers) != 0 {
		app.RegisterMultiple(providers, false)
	}

	app.Boot()

	return app
}

func registerBaseProvider(f contract.Application) {
	f.RegisterMultiple([]contract.Provider{
		new(event.Provider),
		new(logging.Provider),
	}, false)
}
