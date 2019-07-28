package logging

import (
	"github.com/firmeve/firmeve/config"
	"testing"
)

func TestNewLogger(t *testing.T) {
	//sugar := zap.NewExample().Sugar()
	//defer sugar.Sync()
	//sugar.Infow("failed to fetch URL",
	//	"url", "http://example.com",
	//	"attempt", 3,
	//	"backoff", time.Second,
	//)
	//sugar.Infof("failed to fetch URL: %s", "http://example.com")

	//fmt.Println(zapcore.InfoLevel)
	//fmt.Println("================")
	logger := NewLogger(config.NewConfig("../testdata/config"))
	//test := map[string]string{"url":"firmeve.com"}
	logger.Debug("abc", "url", "http://example.com",)
	//logger.Warn("abc")
	logger.Info("abc")
	//logger.Fatal("abc")
	//logger.Error("abc")
}
