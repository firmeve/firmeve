package firmeve

import (
	"github.com/firmeve/firmeve"
	bootstrap2 "github.com/firmeve/firmeve/bootstrap"
	"github.com/firmeve/firmeve/support"
)

func Bootstrap(options ...support.Option) *firmeve.Firmeve {
	f := firmeve.New(firmeve.WithMode(firmeve.ModeTesting))
	bootstrap2.New(f, options...).FastBootFull()
	return f
}
