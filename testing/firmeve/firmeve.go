package firmeve

import (
	"github.com/firmeve/firmeve"
	bootstrap2 "github.com/firmeve/firmeve/bootstrap"
)

func Bootstrap(configPath string) *firmeve.Firmeve {
	f := firmeve.New()

	bootstrap := bootstrap2.New(f, configPath)
	bootstrap.RegisterDefault()
	bootstrap.Boot()

	return f
}
