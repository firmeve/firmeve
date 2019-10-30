package event

import (
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.NewFirmeve()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("event"))
	assert.Equal(t,true,firmeve.Has(`event`))
}
