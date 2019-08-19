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
	Delay(delay time.Time, jobName string, options ...utils.OptionFunc)
	// 出队
	Pop(queueName string) *Payload
}

// 实际队列业务处理者
type Job interface {
	Handle(data interface{}) error
}

// 队列处理器
type Processor interface {
	Handle(payload *Payload)
}

// 队列参数选项
type option struct {
	queueName string
	data      interface{}
	attempt   uint8
	tries     uint8
	timeout   time.Duration
	delay     time.Time
}

// 入队数据载荷
type Payload struct {
	processorName string
	queueName     string
	jobName       string
	timeout       time.Duration
	delay         time.Time
	tries         uint8
	attempt       uint8
	data          interface{}
}

// 队列管理器
type manager struct {
	config      *config.Config
	jobs        map[string]Job
	connections map[string]Queue
	processors  map[string]Processor
}

// 默认处理进程
type processor struct {
}

// Queue error
type Error struct {
	Message string
	payload *Payload
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

func WithAttempt(attempt uint8) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).attempt = attempt
	}
}

func WithDelay(delay time.Time) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).delay = delay
	}
}

func WithTimeout(timeout time.Duration) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).timeout = timeout
	}
}

func WithTries(tries uint8) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).tries = tries
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
			processors:  map[string]Processor{`default`: &processor{}},
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
	mu.Lock()
	if _, ok := m.connections[name]; !ok {
		m.connections[name] = factory(name, m.config)
	}
	mu.Unlock()
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
		attempt:   0,
		delay:     time.Now(),
		timeout:   time.Duration(5) * time.Second,
	}, options...).(*option)

	return &Payload{
		queueName:     option.queueName,
		processorName: ``,
		jobName:       jobName,
		timeout:       option.timeout,
		delay:         option.delay,
		tries:         option.tries,
		attempt:       option.attempt,
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
		payload := m.Connection(`memory`).Pop(queueName)
		if payload != nil {
			queueManager.GetProcess(`default`).Handle(payload)
		}
		//select {
		//case payload := <-m.Connection(`memory`).Pop(queueName):
		//	queueManager.GetProcess(`default`).Handle(payload)
		//case <-time.After(time.Second * 10):
		//	fmt.Println("超时")
		//}
	}
}

// ---------------------------------------------- processor ---------------------------------------

// 默认处理器
func (p *processor) Handle(payload *Payload) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(*Error); !ok {
				panic(err)
			}

			// 设置error信息
			err.(*Error).SetPayload(payload)

			// 重试次数达到，则丢弃否则重新投入队列
			if payload.attempt > payload.tries {
				// panic(`达到最大值`)
				fmt.Println(`队列失败`)
			} else {
				queueManager.Connection(`memory`).Push(payload.jobName, WithAttempt(payload.attempt+1), WithQueueName(payload.queueName), WithData(payload.data))
				fmt.Println("执行", payload.attempt+1)
			}

		}
	}()

	job := queueManager.GetJob(payload.jobName)
	err := job.Handle(payload.data)
	if err != nil {
		panic(err)
	}
}

// ---------------------------------------------- error ---------------------------------------

func (e *Error) SetPayload(payload *Payload) {
	e.payload = payload
}

func (e *Error) Payload() *Payload {
	return e.payload
}

func (e *Error) Error() string {
	return fmt.Sprintf("the queue[%s] execute job[%s] error, message: %s", e.payload.queueName, e.payload.jobName, e.Message)
}
