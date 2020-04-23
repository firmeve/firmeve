package logging

import (
	"fmt"
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
		channels       channels
		config         contract.Configuration
		current        string
		currentChannel internalLogger
	}

	Level string

	internalLogger = *zap.SugaredLogger

	channels map[string]internalLogger

	writers map[string]io.Writer
)

const (
	Debug Level = `debug`
	Info        = `info`
	Warn        = `warn`
	Error       = `error`
	Fatal       = `fatal`
	Panic       = `panic`
)

var (
	levelMap = map[Level]zapcore.Level{
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
	consoleZapEncoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(time.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	})
	fileZapEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
	})
)

func New(config contract.Configuration) contract.Loggable {
	var (
		channels = allChannels(config)
		current  = config.GetString(`default`)
	)
	return &logger{
		config:         config,
		current:        current,
		channels:       channels,
		currentChannel: channels[current],
	}
}

func (l *logger) Debug(message string, context ...interface{}) {
	l.currentChannel.Debugw(message, context...)
}

func (l *logger) Info(message string, context ...interface{}) {
	l.currentChannel.Infow(message, context...)
}

func (l *logger) Warn(message string, context ...interface{}) {
	l.currentChannel.Warnw(message, context...)
}

func (l *logger) Error(message string, context ...interface{}) {
	l.currentChannel.Errorw(message, context...)
}

func (l *logger) Fatal(message string, context ...interface{}) {
	l.currentChannel.Fatalw(message, context...)
}

func (l *logger) Panic(message string, context ...interface{}) {
	l.currentChannel.Panicw(message, context...)
}

// ---------------------------------------------- func --------------------------------------------------

// Default internal logger
func zapLogger(config contract.Configuration, writers writers) internalLogger {
	cores := make([]zapcore.Core, 0)
	var zapEncoder zapcore.Encoder
	for stack, write := range writers {
		if stack == `console` {
			zapEncoder = consoleZapEncoder
		} else {
			zapEncoder = fileZapEncoder
		}

		core := zapcore.NewCore(
			zapEncoder,
			zapcore.Lock(zapcore.AddSync(write)), //writer(option)
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= levelMap[Level(config.GetStringMap(strings.Join([]string{`channels`, stack}, `.`))[`level`].(string))]
			}),
		)

		cores = append(cores, core)
	}

	return zap.New(
		zapcore.NewTee(cores...), zap.AddCallerSkip(3), zap.AddStacktrace(levelMap[Level(config.GetString(`stack_level`))]),
	).Sugar()
}

// Channel factory
func factory(stack string, config contract.Configuration) internalLogger {
	var channels writers
	switch stack {
	case `file`:
		channels = writers{stack: newFileChannel(config)}
	case `console`:
		channels = writers{stack: newConsoleChannel(config)}
	case `stack`:
		channels = newStackChannel(config)
	default:
		panic(fmt.Errorf("the logger stack %s not exists", stack))
	}

	return zapLogger(config, channels)
}

func allChannels(config contract.Configuration) channels {
	stacks := config.Get(`channels`).(map[string]interface{})
	channels := make(channels, 3)

	for key := range stacks {
		channels[key] = factory(key, config)
	}

	return channels
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

// New stack channel
func newStackChannel(config contract.Configuration) writers {
	stacks := config.GetStringSlice(`channels.stack`)
	existsStackMap := make(writers, 0)
	for _, stack := range stacks {
		existsStackMap[stack] = channelMap[stack](config)
	}

	return existsStackMap
}
