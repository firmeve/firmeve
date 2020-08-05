package logging

import (
	config2 "github.com/firmeve/firmeve/config"
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
		app     contract.Application
		logger  *zap.SugaredLogger
		writers writers
		event   contract.Event
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

	event contract.Event
)

func New(app contract.Application) contract.Loggable {
	zapLogger, writers := zapLogger(
		app.Get(`config`).(*config2.Config).Item(`logging`),
	)

	return &logger{
		app:     app,
		logger:  zapLogger,
		writers: writers,
		event:   app.Resolve(`event`).(contract.Event),
	}
}

func (l *logger) Writer(channel string) io.Writer {
	return l.writers[channel]
}

func (l *logger) Logger() *zap.SugaredLogger {
	return l.logger
}

func (l *logger) Debug(values ...interface{}) {
	l.Write(Debug, values...)
}

func (l *logger) Info(values ...interface{}) {
	l.Write(Info, values...)
}

func (l *logger) Warn(values ...interface{}) {
	l.Write(Warn, values...)
}

func (l *logger) Error(values ...interface{}) {
	l.Write(Error, values...)
}

func (l *logger) Fatal(values ...interface{}) {
	l.Write(Fatal, values...)
}

func (l *logger) Panic(values ...interface{}) {
	l.Write(Panic, values...)
}

func (l *logger) With(values ...interface{}) contract.Loggable {
	l.logger = l.logger.With(values...)
	return l
}

func (l *logger) fireLogEvent(level2 level, values []interface{}) {
	if l.event != nil {
		l.event.Dispatch(`logger.before`, l, level2, values)
	}
}

func (l *logger) Write(level2 level, values ...interface{}) {
	l.fireLogEvent(level2, values)
	switch level2 {
	case Debug:
		l.logger.Debug(values...)
	case Info:
		l.logger.Info(values...)
	case Warn:
		l.logger.Warn(values...)
	case Error:
		l.logger.Error(values...)
	case Fatal:
		l.logger.Fatal(values...)
	case Panic:
		l.logger.Panic(values...)
	}
}

// Separated messages
// The first is message
// Other is key => value
func separatedMessages(values ...interface{}) (string, []interface{}) {
	if len(values)/2 == 0 {
		return ``, values
	}

	return values[0].(string), values[1:]
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
