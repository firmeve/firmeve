package bootstrap

import (
	bootstrap2 "github.com/firmeve/firmeve/bootstrap"
	"github.com/firmeve/firmeve/support"
	"github.com/firmeve/firmeve/testing"
)

func Bootstrap(options ...support.Option) *bootstrap2.Bootstrap {
	bootstrap := bootstrap2.New(testing.TestingModeFirmeve(), options...)
	bootstrap.FastBootFull()
	return bootstrap
}
