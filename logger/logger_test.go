package logging

import (
	"fmt"
	"github.com/firmeve/firmeve/container"
	event2 "github.com/firmeve/firmeve/event"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testingConfigPath = "../testdata/config/config.testing.yaml"

var (
	app contract.Application
)

func init() {
	app = kernel.New()
	app.Bind(`application`, app)
	app.Bind(`firmeve`, app)
	app.Bind(`config`, kernel.NewConfig(path.RunRelative(testingConfigPath)), container.WithShare(true))
	//providers = append([]contract.Provider{new(logging.Provider), new(event.Provider)}, providers...)
	app.RegisterMultiple([]contract.Provider{new(event2.Provider), new(Provider)}, false)
	app.Boot()
}

func TestMain(m *testing.M) {
	//app = testing2.ApplicationDefault()
	m.Run()
}

func TestDefault(t *testing.T) {
	logger := Default()
	logger.Debug("Debug", "append", fmt.Sprintf("%v", logger))
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
