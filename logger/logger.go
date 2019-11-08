package logging

import (
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Loggable interface {
	Debug(message string, context ...interface{})
	Info(message string, context ...interface{})
	Warn(message string, context ...interface{})
	Error(message string, context ...interface{})
	Fatal(message string, context ...interface{})
	Channel(stack string) Loggable
}

type logger struct {
	channels channels
	config   *Config
}

type Config struct {
	Current  string
	Channels ConfigChannelType
}

type Level string
type ConfigChannelType map[string]interface{}
type internalLogger = *zap.SugaredLogger
type channels map[string]internalLogger
type writers map[string]io.Writer

const (
	Debug Level = `debug`
	Info        = `info`
	Warn        = `warn`
	Error       = `error`
	Fatal       = `fatal`
)

var (
	mu       sync.Mutex
	levelMap = map[Level]zapcore.Level{
		Debug: zapcore.DebugLevel,
		Info:  zapcore.InfoLevel,
		Warn:  zapcore.WarnLevel,
		Error: zapcore.ErrorLevel,
		Fatal: zapcore.FatalLevel,
	}
	channelMap = map[string]func(config *Config) io.Writer{
		`file`:    newFileChannel,
		`console`: newConsoleChannel,
	}
)

func New(config *Config) Loggable {
	return &logger{
		config:   config,
		channels: make(channels, 0),
	}
}

func Default() Loggable {
	config := &Config{
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
		},
	}

	return New(config)
}

//func (l *logger) Config(config *Config) Loggable {
//	l.config = config
//	return l
//}

func (l *logger) Debug(message string, context ...interface{}) {
	l.channel(l.config.Current).Debugw(message, context...)
}

func (l *logger) Info(message string, context ...interface{}) {
	l.channel(l.config.Current).Infow(message, context...)
}

func (l *logger) Warn(message string, context ...interface{}) {
	l.channel(l.config.Current).Warnw(message, context...)
}

func (l *logger) Error(message string, context ...interface{}) {
	l.channel(l.config.Current).Errorw(message, context...)
}

func (l *logger) Fatal(message string, context ...interface{}) {
	l.channel(l.config.Current).Fatalw(message, context...)
}

// Return a new Logger instance
// But still using internal channels
func (l *logger) Channel(stack string) Loggable {
	return &logger{
		config:   l.config,
		channels: l.channels,
	}
}

// Get designated channel
func (l *logger) channel(stack string) internalLogger {
	if channel, ok := l.channels[stack]; ok {
		return channel
	}

	mu.Lock()
	defer mu.Unlock()

	l.channels[stack] = factory(stack, l.config)
	return l.channels[stack]
}

// ---------------------------------------------- func --------------------------------------------------

// Default internal logger
func zapLogger(config *Config, writers writers) internalLogger {
	//zapcore.EncoderConfig{
	//	TimeKey:        "time",
	//	LevelKey:       "level",
	//	NameKey:        "logger",
	//	CallerKey:      "caller",
	//	MessageKey:     "message",
	//	StacktraceKey:  "stacktrace",
	//	LineEnding:     zapcore.DefaultLineEnding,
	//	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	//	EncodeTime:     zapcore.ISO8601TimeEncoder,
	//	EncodeDuration: zapcore.StringDurationEncoder,
	//	EncodeCaller:   zapcore.FullCallerEncoder,
	//}
	cores := make([]zapcore.Core, 0)
	var zapEncoder zapcore.Encoder
	for stack, write := range writers {
		if stack == `console` {
			zapEncoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		} else {
			zapEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}

		core := zapcore.NewCore(
			zapEncoder,
			zapcore.Lock(zapcore.AddSync(write)), //writer(option)
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= levelMap[Level(config.Channels[stack].(ConfigChannelType)[`level`].(string))]
			}),
		)

		cores = append(cores, core)
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel)).Sugar()
}

// Channel factory
func factory(stack string, config *Config) internalLogger {
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

// New file channel
func newFileChannel(config *Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   config.Channels[`file`].(ConfigChannelType)[`path`].(string) + "/log.log",
		MaxSize:    config.Channels[`file`].(ConfigChannelType)[`size`].(int),
		MaxBackups: config.Channels[`file`].(ConfigChannelType)[`backup`].(int),
		MaxAge:     config.Channels[`file`].(ConfigChannelType)[`age`].(int),
	}
}

// New console channel
func newConsoleChannel(config *Config) io.Writer {
	return os.Stdout
}

// New stack channel
func newStackChannel(config *Config) writers {
	stacks := config.Channels[`stack`].([]string)
	existsStackMap := make(writers, 0)
	for _, stack := range stacks {
		existsStackMap[stack] = channelMap[stack](config)
	}

	return existsStackMap
}
