package testing

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
)

func TestingModeFirmeve() contract.Application {
	app := kernel.New()
	app.SetMode(contract.ModeTesting)
	app.Bind(`firmeve`, app)
	app.Bind(`config`, config.New(path.RunRelative("../testdata/config/config.yaml")))
	return app
}
