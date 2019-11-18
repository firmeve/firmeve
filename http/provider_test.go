package http

import (
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	firmeve := firmeve2.New()
	//firmeve2.BindingInstance(firmeve)
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("http"))
	assert.Equal(t, true, firmeve.Has(`http.router`))
}
