package memory

import (
	"fmt"
	"github.com/firmeve/firmeve/queue"
	"github.com/firmeve/firmeve/utils"
	"time"
)

type Queue string

type Memory struct {
	payload map[string]chan *queue.Payload
}

func NewMemory() *Memory {
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

func (m *Memory) Later(delay time.Time, job queue.Job, queue string, data interface{}) {
}

func (m *Memory) Run(queueName string) {
	processNum := 5

	// 同时开5个进程
	for i := 0; i < processNum; i++ {
		go m.RunProcess(queueName)
	}
}

func (m *Memory) RunProcess(queueName string) {
	for {
		select {
		case _ = <-m.Pop(queueName):
		//case payload := <-m.Pop(queueName):
		// 这块就是worker了
		//job := queue.NewManager().Get(payload.Job)
		//job.Handle(payload.Handle)
		//data.job.Handle(data.data)
		case <-time.After(time.Second * 10):
			fmt.Println("超时")
		}
	}
}
