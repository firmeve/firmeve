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
	Push(job string, options ...utils.OptionFunc)

	//
	PushRaw(payload *Payload)

	// 延时入队
	Delay(delay time.Duration, job string, options ...utils.OptionFunc)

	// 出队
	Pop(queue string) *Payload
}

// 实际队列业务Job
type Job interface {
	Handle(payload *Payload)
	Failed(err *JobError)
}

// 队列处理器
type Processor interface {
	Handle(connection Queue, payload *Payload)
}

// 队列参数选项
type option struct {
	queue     string
	processor string
	data      interface{}
	attempt   uint8
	tries     uint8
	timeout   time.Duration
	timeoutAt time.Time
	delay     time.Time
}

// 入队数据载荷
type Payload struct {
	Processor string
	Queue     string
	Job       string
	Timeout   time.Duration
	TimeoutAt time.Time
	Delay     time.Time
	Tries     uint8
	Attempt   uint8
	Data      interface{}
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
type JobError struct {
	message string
	payload *Payload
	err     error
}

var (
	queueManager *manager
	once         sync.Once
	mu           sync.Mutex
)

func WithQueue(queue string) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).queue = queue
	}
}

func WithProcessor(processor string) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).processor = processor
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

func WithDelay(delay time.Duration) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).delay = time.Now().Add(delay)
	}
}

func WithTimeout(timeout time.Duration) utils.OptionFunc {
	return func(utilOption utils.Option) {
		utilOption.(*option).timeout = timeout
		utilOption.(*option).timeoutAt = time.Now().Add(timeout)
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

// Register handle job
func (m *manager) RegisterJob(name string, job Job) {
	mu.Lock()
	m.jobs[name] = job
	mu.Unlock()
}

// Get handle job
func (m *manager) GetJob(name string) Job {
	return m.jobs[name]
}

// Register handle Process
func (m *manager) RegisterProcess(name string, processor Processor) {
	mu.Lock()
	m.processors[name] = processor
	mu.Unlock()
}

// Get handle Process
func (m *manager) GetProcess(name string) Processor {
	return m.processors[name]
}

// Get a queue connection
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

// Create a core payload
func createPayload(job string, options ...utils.OptionFunc) *Payload {
	option := utils.ApplyOption(&option{
		queue:     `default`,
		processor: `default`,
		data:      nil,
		attempt:   0,
		delay:     time.Now(),
		timeout:   0,
		timeoutAt: time.Now(),
	}, options...).(*option)

	return &Payload{
		Queue:     option.queue,
		Processor: option.processor,
		Job:       job,
		Timeout:   option.timeout,
		TimeoutAt: option.timeoutAt,
		Delay:     option.delay,
		Tries:     option.tries,
		Attempt:   option.attempt,
		Data:      option.data,
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
		// 这里memory后面会用命令行控制
		connection := m.Connection(`memory`)
		payload := connection.Pop(queueName)
		if payload != nil {
			// 这里default后面会用命令行控制
			queueManager.GetProcess(`default`).Handle(connection, payload)
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
func (p *processor) Handle(connection Queue, payload *Payload) {
	job := queueManager.GetJob(payload.Job)

	defer func(job Job) {
		// 再次执行 recover 防止 job.Failed 出现panic
		defer func() {
			if err := recover(); err != nil {
				panic(err)
			}
		}()

		if err := recover(); err != nil {
			var jobError *JobError
			if v, ok := err.(*JobError); !ok {
				jobError = NewJobError()
				if v, ok := err.(error); ok {
					jobError.SetError(v)
				} else if v, ok := err.(string); ok {
					jobError.SetMessage(v)
				}
			} else {
				jobError = v
			}

			jobError.SetPayload(payload)

			// 设置error信息
			job.Failed(jobError)

			// 重试次数达到，则丢弃否则重新投入队列
			if payload.Attempt < payload.Tries {
				// 重新入队
				payload.Attempt = payload.Attempt + 1
				connection.PushRaw(payload)
			}
		}
	}(job)

	// 判断是否超时
	if payload != nil && time.Now().Sub(payload.TimeoutAt).Seconds() < 0 {
		panic(fmt.Sprintf("the job %s timeout", payload.Job))
	}
	// 判断重试次数
	if payload.Tries > 0 && payload.Attempt >= payload.Tries {
		panic(fmt.Sprintf("the job %s execute more than the specified number of times", payload.Job))
	}

	job.Handle(payload)
}

// ---------------------------------------------- error ---------------------------------------
type errorOption struct {
	message string
	err     error
	payload *Payload
}

//func WithJobErrorMessage(message string) utils.OptionFunc {
//	return func(utilOption utils.Option) {
//		utilOption.(*errorOption).message = message
//	}
//}

func NewJobError(options ...utils.OptionFunc) *JobError {
	option := utils.ApplyOption(&errorOption{
		message: "",
		err:     nil,
		payload: nil,
	}, options...).(*errorOption)

	return &JobError{
		message: option.message,
		err:     option.err,
		payload: option.payload,
	}
}

func (e *JobError) SetMessage(message string) *JobError {
	e.message = message
	return e
}

func (e *JobError) SetPayload(payload *Payload) *JobError {
	e.payload = payload
	return e
}

func (e *JobError) SetError(err error) *JobError {
	e.err = err
	return e
}

func (e *JobError) GetMessage() string {
	return e.message
}

func (e *JobError) GetError() error {
	return e.err
}

func (e *JobError) GetPayload() *Payload {
	return e.payload
}

func (e *JobError) Error() string {
	return fmt.Sprintf("the queue[%s] execute job[%s] error, message: %s", e.payload.Queue, e.payload.Job, e.message)
}
