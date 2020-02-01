package testing

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel"
)

func TestingModeFirmeve() kernel.IApplication {
	return firmeve.Default(kernel.ModeTesting, "../testdata/config")
}
