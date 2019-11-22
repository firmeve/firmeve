package bootstrap

import (
	"github.com/firmeve/firmeve"
	bootstrap2 "github.com/firmeve/firmeve/bootstrap"
	"github.com/firmeve/firmeve/support"
	"github.com/firmeve/firmeve/testing"
)

func Bootstrap(options ...support.Option) *firmeve.Firmeve {
	f := testing.TestingModeFirmeve()
	bootstrap2.New(f, options...).FastBootFull()
	return f
}

