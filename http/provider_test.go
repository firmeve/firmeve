package http

import (
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	firmeve := firmeve2.New()
	//firmeve2.BindingInstance(firmeve)
	firmeve.Bind(`config`, config.New(path.RunRelative("../testdata/config")))
	firmeve.Register(firmeve.Make(new(Provider)).(firmeve2.Provider))
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("http"))
	assert.Equal(t, true, firmeve.Has(`http.router`))
}
