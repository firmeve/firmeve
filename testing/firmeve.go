package testing

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
)

var (
	TestingApplication contract.Application
)

func init() {
	TestingApplication = TestingMode()
}

func TestingModeFirmeve() contract.Application {
	return application("../testdata/config/config.testing.yaml")
}

func TestingMode() contract.Application {
	return application("../testdata/config/config.testing.yaml")
}

func TestingModeWithConfig(configPath string) contract.Application {
	return application(configPath)
}

func application(configPath string) contract.Application {
	app := kernel.New()
	app.SetMode(contract.ModeTesting)
	app.Bind(`firmeve`, app)
	app.Bind(`config`, config.New(path.RunRelative(configPath)))
	return app
}
