package logging

import (
	"github.com/firmeve/firmeve/config"
	"testing"
)

func TestNewLogger(t *testing.T) {
	//fmt.Println(zapcore.InfoLevel)
	//fmt.Println("================")
	logger := NewLogger(config.NewConfig("../testdata/config"))
	logger.Debug("abc")
	logger.Warn("abc")
	logger.Info("abc")
	logger.Fatal("abc")
	logger.Error("abc")
}
