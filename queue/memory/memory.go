package memory

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/queue"
	"github.com/firmeve/firmeve/utils"
	"time"
)


type Memory struct {
	payload map[string]chan *queue.Payload
}

func NewMemory(config *config.Config) *Memory {
	return &Memory{
		payload: make(map[string]chan *queue.Payload),
	}
}

func (m *Memory) Push(jobName string, options ...utils.OptionFunc) {
	//
	option := utils.ApplyOption(&queue.Option{
		QueueName: `default`,
		Data:      nil,
	}, options...).(*queue.Option)

	if _, ok := m.payload[option.QueueName]; !ok {
		m.payload[option.QueueName] = make(chan *queue.Payload)
	}

	m.payload[option.QueueName] <- &queue.Payload{
		JobName:  jobName,
		Timeout:  time.Now(),
		Delay:    time.Now(),
		MaxTries: 8,
		Data:     option.Data,
	}
}

func (m *Memory) Pop(queueName string) <-chan *queue.Payload {
	return m.payload[queueName]
}

func (m *Memory) Later(delay time.Time, jobName string, options ...utils.OptionFunc) {
}
