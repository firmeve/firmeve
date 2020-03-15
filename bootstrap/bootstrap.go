package bootstrap

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/kernel/contract"
	logging "github.com/firmeve/firmeve/logger"
)

func Boot(command contract.Command) contract.Application {
	cmd := command.Cmd()
	app := command.Application()
	providers := command.Providers()

	configPath := cmd.Flag(`config`).Value.String()
	devMode := cmd.Flag(`dev`).Value.String()
	var mode uint8
	if devMode == `true` {
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
