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
	manager  *Manager
	once     sync.Once
)

func NewLogger(config *config.Config) *Manager {
	if manager != nil {
		return manager
	}

	once.Do(func() {
		manager = &Manager{
			channels: resolveChannels(config.Item(`logging`)),
		}
	})

	return manager
}

func factory(stack string, config *config.Config) Logger {
	var logger Logger
	switch stack {
	case `file`:
		logger = newFileLogger(
			Debug,
			&fileOption{
				file:   config.GetString("channels.file.path") + ".log",
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

func resolveChannels(config *config.Config) Channel {

	stacks := config.Item(`logging`).GetStringSlice(`stacks`)
	stackChannels := make(Channel, len(stacks))
	for _, stack := range stacks {
		stackChannels[stack] = factory(stack, config)
	}

	return stackChannels
}

func (m *Manager) Debug(message string, context ...interface{}) {
	for _, channel := range m.channels {
		channel.Debug(message, context)
	}
}

func (m *Manager) Info(message string, context ...interface{}) {
	for _, channel := range m.channels {
		channel.Debug(message, context)
	}
}

func (m *Manager) Warn(message string, context ...interface{}) {
	for _, channel := range m.channels {
		channel.Debug(message, context)
	}
}

func (m *Manager) Error(message string, context ...interface{}) {
	for _, channel := range m.channels {
		channel.Debug(message, context)
	}
}

func (m *Manager) Fatal(message string, context ...interface{}) {
	for _, channel := range m.channels {
		channel.Debug(message, context)
	}
}

//func (m *Manager) newChannel(stack string) Logger {
//	//return factory(stack,)
//}

func (m *Manager) channel(stack string) Logger {
	return m.channels[stack]
}