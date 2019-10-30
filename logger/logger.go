package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

type LoggerInterface interface {
	Debug(message string, context ...interface{})
	Info(message string, context ...interface{})
	Warn(message string, context ...interface{})
	Error(message string, context ...interface{})
	Fatal(message string, context ...interface{})
}

type Logger struct {
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
	logger  *Logger
	mu       sync.Mutex
	once     sync.Once
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

func New(config *Config) *Logger {
	if logger != nil {
		return logger
	}

	once.Do(func() {
		logger = &Logger{
			config:   config,
			channels: make(channels, 0),
		}
	})

	return logger
}

func Default() *Logger {
	config := &Config{
		Current: `console`,
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

func (l *Logger) Config(config *Config) *Logger {
	l.config = config
	return l
}

func (l *Logger) Debug(message string, context ...interface{}) {
	l.channel(l.config.Current).Debugw(message, context...)
}

func (l *Logger) Info(message string, context ...interface{}) {
	l.channel(l.config.Current).Infow(message, context...)
}

func (l *Logger) Warn(message string, context ...interface{}) {
	l.channel(l.config.Current).Warnw(message, context...)
}

func (l *Logger) Error(message string, context ...interface{}) {
	l.channel(l.config.Current).Errorw(message, context...)
}

func (l *Logger) Fatal(message string, context ...interface{}) {
	l.channel(l.config.Current).Fatalw(message, context...)
}

// Get designated channel
func (l *Logger) channel(stack string) internalLogger {
	if channel, ok := l.channels[stack]; ok {
		return channel
	}

	mu.Lock()
	defer mu.Unlock()

	l.channels[stack] = factory(stack, l.config)
	return l.channels[stack]
}

// Return a new Logger instance
// But still using internal channels
func (l *Logger) Channel(stack string) LoggerInterface {
	return &Logger{
		config:   l.config,
		channels: l.channels,
	}
}

// ---------------------------------------------- func --------------------------------------------------

// Default internal logger
func zapLogger(config *Config, writers writers) internalLogger {
	cores := make([]zapcore.Core, 0)
	for stack, write := range writers {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.FullCallerEncoder,
			}),
			zapcore.Lock(zapcore.AddSync(write)), //writer(option)
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= levelMap[Level(config.Channels[stack].(ConfigChannelType)[`level`].(string))]
			}),
		)

		cores = append(cores, core)
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel)).Sugar()
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
