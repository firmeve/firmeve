package http

import (
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
	"testing"
)

var configPath = "../testdata/config"

func TestProvider(t *testing.T) {
	firmeve := firmeve2.New(kernel.ModeProduction, configPath)
	//firmeve2.BindingInstance(firmeve)
	firmeve.Bind(`config`, config.New(path.RunRelative("../testdata/config")))
	firmeve.Register(firmeve.Make(new(Provider)).(kernel.IProvider), false)
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("http"))
	assert.Equal(t, true, firmeve.Has(`http.router`))
}
