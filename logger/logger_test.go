package logging

import (
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestLogger_Channel(t *testing.T) {
	Default().Channel(`file`).Debug("Debug")
}

func TestLogger_Logger_Config(t *testing.T) {
	logger := New(&Config{
		Current: `stack`,
		Channels: ConfigChannelType{
			`stack`: []string{`file`, `console`},
			`console`: ConfigChannelType{
				`level`: `debug`,
			},
			`file`: ConfigChannelType{
				`level`:  `debug`,
				`path`:   "../testdata/logs",
				`size`:   100,
				`backup`: 3,
				`age`:    1,
			},
		}})

	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")
	//logger.Fatal("Fatal")

	assert.Equal(t, true, true)
}

func TestLogger_File(t *testing.T) {
	logger := New(&Config{
		Current: `file`,
		Channels: ConfigChannelType{
			`stack`: []string{`file`, `console`},
			`console`: ConfigChannelType{
				`level`: `debug`,
			},
			`file`: ConfigChannelType{
				`level`:  `debug`,
				`path`:   "../testdata/logs",
				`size`:   100,
				`backup`: 3,
				`age`:    1,
			},
		}})

	logger.Warn("File")
}


func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.Instance()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("logger"))
	assert.Equal(t,true,firmeve.Has(`logger`))
}