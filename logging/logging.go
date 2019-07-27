package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	zap *zap.Logger
}

func NewLogger() {
	//core := zapcore.NewTee(
	//	//zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
	//	//zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
	//	//zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
	//	//zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	//)

	z := zapcore.NewTee(
		file("a"),
		file("b"),
	)
	logger := zap.New(z)
	defer logger.Sync()
	logger.Info("constructed a logger")
}

func file(name string) zapcore.Core {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app" +name+".log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     1,
	})
	config := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	return zapcore.NewCore(config,writer,lowPriority)
}
