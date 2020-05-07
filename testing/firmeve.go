package testing

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
)

const testingConfigPath = "../testdata/config/config.testing.yaml"

//
//var (
//	TestingApplication contract.Application
//)
//
//func init() {
//	TestingApplication = TestingMode()
//}
//
//func TestingModeFirmeve() contract.Application {
//	return application("../testdata/config/config.testing.yaml")
//}
//
//func TestingMode() contract.Application {
//	return application("../testdata/config/config.testing.yaml")
//}
//
//func TestingModeWithConfig(configPath string) contract.Application {
//	return application(configPath)
//}
//

func Application(configPath string, providers ...contract.Provider) contract.Application {
	app := kernel.New()
	bootstrap(app, configPath, providers...)
	return app
}

func ApplicationDefault(providers ...contract.Provider) contract.Application {
	return Application(testingConfigPath, providers...)
}

func bootstrap(app contract.Application, configPath string, providers ...contract.Provider) {
	app.SetMode(contract.ModeTesting)
	app.Bind(`firmeve`, app)
	app.Bind(`config`, config.New(configPath), container.WithShare(true))
	app.RegisterMultiple(providers, false)
	app.Boot()
}

//
//func application(configPath string) contract.Application {
//	app := kernel.New()
//	app.SetMode(contract.ModeTesting)
//	app.Bind(`firmeve`, app)
//	app.Bind(`config`, config.New(path.RunRelative(configPath)))
//	return app
//}
