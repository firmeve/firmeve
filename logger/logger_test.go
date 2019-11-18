package logging

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/support/path"
	"testing"

	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	logger := Default()
	logger.Debug("Debug")
	logger.Info("Info")
	logger.Error("Error")
	logger.Warn("Warn")
	//
	//assert.Equal(t, true, true)
}

func Default() Loggable {
	return New(config.New(path.RunRelative("../testdata/config")).Item(`logging`))
}

func TestLogger_Channel(t *testing.T) {
	Default().Channel(`file`).Debug("Debug")
}

func TestLogger_Logger_Config(t *testing.T) {
	logger := New(config.New(path.RunRelative("../testdata/config")).Item(`logging`))

	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")
	//logger.Fatal("Fatal")

	assert.Equal(t, true, true)
}

func TestLogger_File(t *testing.T) {
	logger := New(config.New(path.RunRelative("../testdata/config")).Item(`logging`))

	logger.Warn("File")
}

func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.New()
	firmeve.Bind(`config`, config.New(path.RunRelative("../testdata/config")))
	firmeve.Register(firmeve.Make(new(Provider)).(firmeve2.Provider))
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("logger"))
	assert.Equal(t, true, firmeve.Has(`logger`))
}
