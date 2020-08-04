package logging

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

type (
	logger struct {
		writers writers
		logger  *zap.SugaredLogger
	}

	level string

	writers map[string]io.Writer
)

const (
	Debug level = `debug`
	Info        = `info`
	Warn        = `warn`
	Error       = `error`
	Fatal       = `fatal`
	Panic       = `panic`
)

var (
	levelMap = map[level]zapcore.Level{
		Debug: zapcore.DebugLevel,
		Info:  zapcore.InfoLevel,
		Warn:  zapcore.WarnLevel,
		Error: zapcore.ErrorLevel,
		Fatal: zapcore.FatalLevel,
		Panic: zapcore.PanicLevel,
	}
	channelMap = map[string]func(config contract.Configuration) io.Writer{
		`file`:    newFileChannel,
		`console`: newConsoleChannel,
	}
	zapEncoderConfig = zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(time.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	zapEncoders = map[string]zapcore.Encoder{
		`console`: zapcore.NewConsoleEncoder(zapEncoderConfig),
		`json`:    zapcore.NewJSONEncoder(zapEncoderConfig),
	}
)

func New(config contract.Configuration) contract.Loggable {
	zapLogger, writers := zapLogger(config)

	return &logger{
		logger:  zapLogger,
		writers: writers,
	}
}

func (l *logger) Writer(channel string) io.Writer {
	return l.writers[channel]
}

func (l *logger) Logger() *zap.SugaredLogger {
	return l.logger
}

func (l *logger) Debug(message string, context ...interface{}) {
	l.logger.Debugw(message, context...)
}

func (l *logger) Info(message string, context ...interface{}) {
	l.logger.Infow(message, context...)
}

func (l *logger) Warn(message string, context ...interface{}) {
	l.logger.Warnw(message, context...)
}

func (l *logger) Error(message string, context ...interface{}) {
	l.logger.Errorw(message, context...)
}

func (l *logger) Fatal(message string, context ...interface{}) {
	l.logger.Fatalw(message, context...)
}

func (l *logger) Panic(message string, context ...interface{}) {
	l.logger.Panicw(message, context...)
}

func (l *logger) With(context ...interface{}) contract.Loggable {
	l.logger = l.logger.With(context...)
	return l
}

// ---------------------------------------------- func --------------------------------------------------

// Default internal logger
func zapLogger(config contract.Configuration) (*zap.SugaredLogger, writers) {
	encoder := zapEncoders[config.GetString(`formatter`)]
	current := config.GetStringSlice(`default`)
	writers := make(writers, len(current))
	cores := make([]zapcore.Core, len(current))
	for i := range current {
		writers[current[i]] = channelMap[current[i]](config)
		cores[i] = zapcore.NewCore(
			encoder,
			zapcore.AddSync(writers[current[i]]),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= levelMap[level(config.GetStringMap(strings.Join([]string{`channels`, current[i]}, `.`))[`level`].(string))]
			}))
	}

	return zap.New(
		zapcore.NewTee(cores...), zap.AddCallerSkip(3), zap.AddStacktrace(levelMap[level(config.GetString(`stack_level`))]),
	).Sugar(), writers
}

// New file channel
func newFileChannel(config contract.Configuration) io.Writer {
	return &lumberjack.Logger{
		Filename:   config.GetStringMap(`channels.file`)[`path`].(string) + "/log.log",
		MaxSize:    config.GetStringMap(`channels.file`)[`size`].(int),
		MaxBackups: config.GetStringMap(`channels.file`)[`backup`].(int),
		MaxAge:     config.GetStringMap(`channels.file`)[`age`].(int),
	}
}

// New console channel
func newConsoleChannel(config contract.Configuration) io.Writer {
	return os.Stdout
}
