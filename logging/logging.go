package logging

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

type Logger interface {
	Debug(message string, context ...interface{})
	Info(message string, context ...interface{})
	Warn(message string, context ...interface{})
	Error(message string, context ...interface{})
	Fatal(message string, context ...interface{})
}

type Manager struct {
	channels channels
	current string
	config  *config.Config
}

type Level string
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
	manager  *Manager
	mu       sync.Mutex
	once     sync.Once
	levelMap = map[Level]zapcore.Level{
		Debug: zapcore.DebugLevel,
		Info:  zapcore.InfoLevel,
		Warn:  zapcore.WarnLevel,
		Error: zapcore.ErrorLevel,
		Fatal: zapcore.FatalLevel,
	}
	channelMap = map[string]func(config *config.Config) io.Writer{
		`file`:    newFileChannel,
		`console`: newConsoleChannel,
	}
)

func NewLogger(config *config.Config) *Manager {
	if manager != nil {
		return manager
	}

	config = config.Item(`logging`)

	once.Do(func() {
		manager = &Manager{
			config:   config,
			channels: make(channels, 0),
			current:  config.GetString(`default`),
		}
	})

	return manager
}

func (m *Manager) Debug(message string, context ...interface{}) {
	m.channel(m.current).Debugw(message, context...)
}

func (m *Manager) Info(message string, context ...interface{}) {
	m.channel(m.current).Infow(message, context...)
}

func (m *Manager) Warn(message string, context ...interface{}) {
	m.channel(m.current).Warnw(message, context...)
}

func (m *Manager) Error(message string, context ...interface{}) {
	m.channel(m.current).Errorw(message, context...)
}

func (m *Manager) Fatal(message string, context ...interface{}) {
	m.channel(m.current).Fatalw(message, context...)
}

// Get designated channel
func (m *Manager) channel(stack string) internalLogger {
	if channel, ok := m.channels[stack]; ok {
		return channel
	}

	mu.Lock()
	defer mu.Unlock()

	m.channels[stack] = factory(stack, m.config)
	return m.channels[stack]
}

// Return a new Logger instance
// But still using internal channels
func (m *Manager) Channel(stack string) Logger {
	return &Manager{
		config:   m.config,
		channels: m.channels,
		current:  stack,
	}
}

// ---------------------------------------------- func --------------------------------------------------

// Default internal logger
func zapLogger(config *config.Config, writers writers) internalLogger {
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
				return lvl >= levelMap[Level(config.GetString(`channels.`+stack+`.level`))]
			}),
		)

		cores = append(cores, core)
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel)).Sugar()
}

// Channel factory
func factory(stack string, config *config.Config) internalLogger {
	var channels writers
	switch stack {
	case `file`:
		channels = writers{stack: newFileChannel(config)}
	case `console`:
		channels = writers{stack: newConsoleChannel(config)}
	case `stack`:
		channels = newStackChannel(config)
	default:
		panic(fmt.Sprintf("the logger stack %s not exists", stack))
	}

	return zapLogger(config, channels)
}

// New file channel
func newFileChannel(config *config.Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   config.GetString("channels.file.path") + "/log.log",
		MaxSize:    config.GetInt("channels.file.size"),
		MaxBackups: config.GetInt("channels.file.backup"),
		MaxAge:     config.GetInt("channels.file.age"),
	}
}

// New console channel
func newConsoleChannel(config *config.Config) io.Writer {
	return os.Stdout
}

// New stack channel
func newStackChannel(config *config.Config) writers {
	stacks := config.GetStringSlice("channels.stack")
	existsStackMap := make(writers, 0)
	for _, stack := range stacks {
		existsStackMap[stack] = channelMap[stack](config)
	}

	return existsStackMap
}
