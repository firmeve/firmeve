package queue

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/utils"
	"time"
)

type memory struct {
	list map[string]chan *Payload
}

func NewMemory(config *config.Config) *memory {
	return &memory{
		list: make(map[string]chan *Payload),
	}
}

func (m *memory) Push(jobName string, options ...utils.OptionFunc) {
	payloadBlock := createPayload(jobName, options...)

	if _, ok := m.list[payloadBlock.queueName]; !ok {
		mu.Lock()
		m.list[payloadBlock.queueName] = make(chan *Payload)
		mu.Unlock()
	}

	m.list[payloadBlock.queueName] <- payloadBlock
}

func (m *memory) Pop(queueName string) <-chan *Payload {
	return m.list[queueName]
}

func (m *memory) Later(delay time.Time, jobName string, options ...utils.OptionFunc) {
}
