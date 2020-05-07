package testing

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
)

const testingConfigPath = "../testdata/config/config.testing.yaml"

func Application(configPath string, providers ...contract.Provider) contract.Application {
	app := kernel.New()
	bootstrap(app, configPath, providers...)
	return app
}

func ApplicationDefault(providers ...contract.Provider) contract.Application {
	return Application(path.RunRelative(testingConfigPath), providers...)
}

func bootstrap(app contract.Application, configPath string, providers ...contract.Provider) {
	app.SetMode(contract.ModeTesting)
	app.Bind(`firmeve`, app)
	app.Bind(`config`, config.New(configPath), container.WithShare(true))
	app.RegisterMultiple(providers, false)
	app.Boot()
}
