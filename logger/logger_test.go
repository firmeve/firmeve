package logging

import (
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	app contract.Application
)

func TestMain(m *testing.M) {
	app = testing2.ApplicationDefault(new(Provider))
	m.Run()
}

func TestDefault(t *testing.T) {
	logger := Default()
	logger.Debug("Debug")
	logger.Info("Info")
	logger.Error("Error")
	logger.Warn("Warn")
	assert.Panics(t, func() {
		logger.Panic("panic")
	})
	//assert.Fail(t,"")
	//
	//assert.Equal(t, true, true)
}

func Default() contract.Loggable {
	return app.Resolve(`logger`).(contract.Loggable)
}

//func TestLogger_Channel(t *testing.T) {
//	Default().Channel(`file`).Debug("Debug")
//}

func TestLogger_Logger_Config(t *testing.T) {
	logger := Default()

	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")
	//logger.Fatal("Fatal")

	assert.Equal(t, true, true)
}

func TestLogger_File(t *testing.T) {
	logger := Default()

	logger.Warn("File")
}
