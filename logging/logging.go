package logging

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
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
	channels Channel
	//default channel
	channel string
	config  *config.Config
}

type Level int8
type Channel map[string]Logger

const (
	Debug Level = iota
	Info        = iota
	Warn        = iota
	Error       = iota
	Fatal       = iota
)

var (
	manager *Manager
	mu      sync.Mutex
	once    sync.Once
)

func NewLogger(config *config.Config) *Manager {
	if manager != nil {
		return manager
	}

	once.Do(func() {
		manager = &Manager{
			config:   config.Item(`logging`),
			channel:  config.Item(`logging`).GetString(`default`),
			channels: make(Channel, 0),
		}
	})

	return manager
}

func (m *Manager) Debug(message string, context ...interface{}) {
	m.Channel(m.channel).Debug(message, context...)
}

func (m *Manager) Info(message string, context ...interface{}) {
	m.Channel(m.channel).Debug(message, context...)
}

func (m *Manager) Warn(message string, context ...interface{}) {
	m.Channel(m.channel).Debug(message, context...)
}

func (m *Manager) Error(message string, context ...interface{}) {
	m.Channel(m.channel).Debug(message, context...)
}

func (m *Manager) Fatal(message string, context ...interface{}) {
	m.Channel(m.channel).Debug(message, context...)
}

func (m *Manager) Channel(stack string) Logger {
	if channel, ok := m.channels[stack]; ok {
		return channel
	}

	mu.Lock()
	defer mu.Unlock()

	m.channels[stack] = factory(stack, m.config)
	return m.channels[stack]
}

// ---------------------------------------------- func --------------------------------------------------

func factory(stack string, config *config.Config) Logger {
	var logger Logger
	switch stack {
	case `file`:
		fmt.Println(config.GetString("channels.file.path") + "/log.log")
		logger = newFileLogger(
			Debug,
			&fileOption{
				file:   config.GetString("channels.file.path") + "/log.log",
				size:   config.GetInt("channels.file.size"),
				backup: config.GetInt("channels.file.backup"),
				age:    config.GetInt("channels.file.age"),
			},
		)
	default:
		panic(fmt.Sprintf("the logger stack %s not exists", stack))
	}

	return logger
}
