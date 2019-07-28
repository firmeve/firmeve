package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type file struct {
	zap *zap.Logger
}

type fileOption struct {
	file   string
	size   int
	backup int
	age    int
	// @todo 后期增加
	sugar  bool
}

var (
	loggerMap = map[Level]zapcore.Level{
		Debug: zapcore.DebugLevel,
		Info:  zapcore.InfoLevel,
		Warn:  zapcore.WarnLevel,
		Error: zapcore.ErrorLevel,
		Fatal: zapcore.FatalLevel,
	}
)

func newFileLogger(level Level, option *fileOption) Logger {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   option.file,
		MaxSize:    option.size,
		MaxBackups: option.backup,
		MaxAge:     option.age,
	})
	//zap.NewProductionEncoderConfig()
	//zapcore.Dev

	core := zapcore.NewCore(
		// 暂时使用production
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}),
		writer,
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == loggerMap[level]
		}),
	)

	return &file{
		zap: zap.New(core,zap.Development(),zap.AddCaller(),zap.AddStacktrace(zap.DebugLevel)),
	}
}

func (f *file) Debug(message string, context ...interface{}) {
	f.zap.Sugar().Debugw(message, context...)
	f.zap.Sync()
}

func (f *file) Info(message string, context ...interface{}) {
	f.zap.Sugar().Infow(message, context...)
	f.zap.Sync()
}

func (f *file) Warn(message string, context ...interface{}) {
	f.zap.Sugar().Warnw(message, context...)
	f.zap.Sync()
}

func (f *file) Error(message string, context ...interface{}) {
	f.zap.Sugar().Errorw(message, context...)
	f.zap.Sync()
}

func (f *file) Fatal(message string, context ...interface{}) {
	f.zap.Sugar().Fatalw(message, context...)
	f.zap.Sync()
}