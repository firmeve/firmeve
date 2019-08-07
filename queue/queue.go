package queue

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/utils"
	"sync"
	"time"
)

type Queue interface {
	// 入队
	Push(jobName string, options ...utils.OptionFunc)
	// 延时入队
	Later(delay time.Time, jobName string, options ...utils.OptionFunc)
	// 出队
	Pop(queueName string) <-chan *Payload
}

// 实际队列业务处理者
type Job interface {
	Handle(data interface{})
}

// 队列处理器
type Processor interface {
	Handle(payload *Payload)
}

// 队列参数选项
type option struct {
	queueName string
	data      interface{}
}

// 入队数据载荷
type Payload struct {
	processorName string
	queueName     string
	jobName       string
	timeout       time.Time
	delay         time.Time
	maxTries      uint8
	data          interface{}
}

// 队列管理器
type manager struct {
	config      *config.Config
	jobs        map[string]Job
	connections map[string]Queue
	processors  map[string]Processor
}

type processor struct {

}

var (
	queueManager *manager
	once         sync.Once
	mu           sync.Mutex
)

func WithQueueName(queueName string) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).queueName = queueName
	}
}

func WithData(data interface{}) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).data = data
	}
}

func NewManager(config *config.Config) *manager {
	if queueManager != nil {
		return queueManager
	}

	once.Do(func() {
		queueManager = &manager{
			config:      config.Item(`queue`),
			jobs:        make(map[string]Job),
			connections: make(map[string]Queue),
			processors:  map[string]Processor{`default`:&processor{}},
		}
	})

	return queueManager
}

func (m *manager) RegisterJob(name string, job Job) {
	mu.Lock()
	m.jobs[name] = job
	mu.Unlock()
}

func (m *manager) GetJob(name string) Job {
	return m.jobs[name]
}

func (m *manager) RegisterProcess(name string, processor Processor) {
	mu.Lock()
	m.processors[name] = processor
	mu.Unlock()
}

func (m *manager) GetProcess(name string) Processor {
	return m.processors[name]
}

func (m *manager) Connection(name string) Queue {
	if _, ok := m.connections[name]; !ok {
		m.connections[name] = factory(name, m.config)
	}

	return m.connections[name]
}

// Connection factory
func factory(name string, config *config.Config) Queue {
	switch name {
	case `memory`:
		return NewMemory(config)
	default:
		panic(fmt.Sprintf("the queue driver %s not exists", name))
	}
}

func createPayload(jobName string, options ...utils.OptionFunc) *Payload {
	option := utils.ApplyOption(&option{
		queueName: `default`,
		data:      nil,
	}, options...).(*option)

	return &Payload{
		queueName:     option.queueName,
		processorName: ``,
		jobName:       jobName,
		timeout:       time.Now(),
		delay:         time.Now(),
		maxTries:      8,
		data:          option.data,
	}
}

func (m *manager) Run(queueName string) {
	processNum := 5

	// 同时开5个进程
	for i := 0; i < processNum; i++ {
		go m.RunProcess(queueName)
	}
}

func (m *manager) RunProcess(queueName string) {
	for {
		select {
		case payload := <-m.Connection(`memory`).Pop(queueName):
			// recover(), error() 如果出错，并且有重试，重新放回队列

			queueManager.GetProcess(`default`).Handle(payload)
			//NewWorker(payload).Handle()
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

func (p *processor) Handle(payload *Payload) {
	job := queueManager.GetJob(payload.jobName)
	job.Handle(payload.data)
}
