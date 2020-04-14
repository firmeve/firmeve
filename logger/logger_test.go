package logging

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	configPath = "../testdata/config/config.yaml"
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

func Default() contract.Loggable {
	return New(config.New(path.RunRelative(configPath)).Item(`logging`))
}

func TestLogger_Channel(t *testing.T) {
	Default().Channel(`file`).Debug("Debug")
}

func TestLogger_Logger_Config(t *testing.T) {
	logger := New(config.New(path.RunRelative(configPath)).Item(`logging`))

	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")
	//logger.Fatal("Fatal")

	assert.Equal(t, true, true)
}

func TestLogger_File(t *testing.T) {
	logger := New(config.New(path.RunRelative(configPath)).Item(`logging`))

	logger.Warn("File")
}

//func TestProvider_Register(t *testing.T) {
//	firmeve := testing2.TestingModeFirmeve()
//	firmeve.Register(new(Provider),true)
//	firmeve.Boot()
//	assert.Equal(t, true, firmeve.HasProvider("logger"))
//	assert.Equal(t, true, firmeve.Has(`logger`))
//}
