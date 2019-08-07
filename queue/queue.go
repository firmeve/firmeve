package queue

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/queue/memory"
	"github.com/firmeve/firmeve/utils"
	"sync"
	"time"
)

type Queue interface {
	Push(jobName string, options ...utils.OptionFunc)
	Pop(queueName string) <-chan *Payload
	Later(delay time.Time, jobName string, options ...utils.OptionFunc)
}

type Option struct {
	QueueName string
	Data      interface{}
}

type Payloader map[string]chan *Payload

type Job interface {
	Handle(data interface{})
}

type Payload struct {
	QueueName string
	JobName  string
	Timeout  time.Time
	Delay    time.Time
	MaxTries uint8
	Data     interface{}
}

var (
	manager *Manager
	once    sync.Once
	jobs    map[string]Job = make(map[string]Job)
)

func WithQueueName(queueName string) utils.OptionFunc {
	return func(option utils.Option) {
		option.(*Option).QueueName = queueName
	}
}

func WithData(data interface{}) utils.OptionFunc {
	return func(option utils.Option) {
		option.(*Option).Data = data
	}
}

type Manager struct {
	config      *config.Config
	//jobs        map[string]Job
	connections map[string]Queue
}

func NewManager(config *config.Config) *Manager {
	if manager != nil {
		return manager
	}

	once.Do(func() {
		manager = &Manager{
			config:      config.Item(`queue`),
			//jobs:        make(map[string]Job),
			connections: make(map[string]Queue),
		}
	})

	return manager
}

func RegisterJob(jobName string, job Job) {
	jobs[jobName] = job
}

func GetJob(jobName string) Job {
	return jobs[jobName]
}

func (m *Manager) Connection(name string) Queue {
	if _, ok := m.connections[name]; !ok {
		m.connections[name] = factory(name, m.config)
	}

	return m.connections[name]
}

func factory(name string, config *config.Config) Queue {
	switch name {
	case `memory`:
		return memory.NewMemory(config)
	default:
		panic(fmt.Sprintf("the queue driver %s not exists", name))
	}
}

func NewPayload(jobName string, options ...utils.OptionFunc) *Payload  {
	option := utils.ApplyOption(&Option{
		QueueName: `default`,
		Data:      nil,
	}, options...).(*Option)

	return &Payload{
		QueueName: option.QueueName,
		JobName:  jobName,
		Timeout:  time.Now(),
		Delay:    time.Now(),
		MaxTries: 8,
		Data:     option.Data,
	}
}


func Run(queueName string) {
	processNum := 5

	// 同时开5个进程
	for i := 0; i < processNum; i++ {
		go RunProcess(queueName)
	}
}

func RunProcess(queueName string) {
	for {
		select {
		case payload := <-manager.Connection(`memory`).Pop(queueName):
			// recover(), error() 如果出错，并且有重试，重新放回队列

			NewWorker(payload).Handle()
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
