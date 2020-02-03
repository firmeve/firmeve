package testing

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/support/path"
)

func TestingModeFirmeve() kernel.IApplication {
	app := kernel.New(kernel.ModeTesting)
	app.Bind(`firmeve`,app)
	app.Bind(`config`, config.New(path.RunRelative("../testdata/config")))
	return app
}
