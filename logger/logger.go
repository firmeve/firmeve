package logging

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
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

	Configuration struct {
		Default    []string               `json:"default" yaml:"default"`
		Channels   map[string]interface{} `json:"channels" yaml:"channels"`
		StackLevel string                 `json:"stack_level" yaml:"stack_level" mapstructure:"stack_level"`
		Formatter  string                 `json:"formatter" yaml:"formatter"`
	}
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
	channelMap = map[string]func(config *Configuration) io.Writer{
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

func New(config *Configuration, event contract.Event) contract.Loggable {
	zapLogger, writers := zapLogger(config)

	return &logger{
		//app:     app,
		logger:  zapLogger,
		writers: writers,
		event:   event,
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

func (l *logger) fireLogEvent(level2 level, message string, values []interface{}) {
	if l.event != nil {
		l.event.Dispatch(`logger.before`, l, level2, message, values)
	}
}

func (l *logger) Write(level2 level, values ...interface{}) {
	message, values := separatedMessages(values...)
	l.fireLogEvent(level2, message, values)
	switch level2 {
	case Debug:
		l.logger.Debugw(message, values...)
	case Info:
		l.logger.Infow(message, values...)
	case Warn:
		l.logger.Warnw(message, values...)
	case Error:
		l.logger.Errorw(message, values...)
	case Fatal:
		l.logger.Fatalw(message, values...)
	case Panic:
		l.logger.Panicw(message, values...)
	}
}

// Separated messages
// The first is message
// Other is key => value
func separatedMessages(values ...interface{}) (string, []interface{}) {
	if len(values)%2 == 0 {
		return ``, values
	}

	return values[0].(string), values[1:]
}

// ---------------------------------------------- func --------------------------------------------------

// Default internal logger
func zapLogger(config *Configuration) (*zap.SugaredLogger, writers) {
	encoder := zapEncoders[config.Formatter]
	current := config.Default
	writers := make(writers, len(current))
	cores := make([]zapcore.Core, len(current))
	for i := range current {
		writers[current[i]] = channelMap[current[i]](config)
		cores[i] = zapcore.NewCore(
			encoder,
			zapcore.AddSync(writers[current[i]]),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= levelMap[level(config.Channels[current[i]].(map[string]interface{})[`level`].(string))]
			}))
	}

	return zap.New(
		zapcore.NewTee(cores...), zap.AddCallerSkip(3), zap.AddStacktrace(levelMap[level(config.StackLevel)]),
	).Sugar(), writers
}

// New file channel
func newFileChannel(config *Configuration) io.Writer {
	fileConfig := config.Channels[`file`].(map[string]interface{})
	return &lumberjack.Logger{
		Filename:   fileConfig[`path`].(string) + "/log.log",
		MaxSize:    fileConfig[`size`].(int),
		MaxBackups: fileConfig[`backup`].(int),
		MaxAge:     fileConfig[`age`].(int),
	}
}

// New console channel
func newConsoleChannel(config *Configuration) io.Writer {
	return os.Stdout
}
