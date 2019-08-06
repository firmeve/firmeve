package queue

import (
	"github.com/firmeve/firmeve/utils"
	"image/gif"
	"time"
)

type Queue interface {
	Push(jobName string, options ...utils.OptionFunc)
	Pop(queueName string)
	Later(delay time.Time, jobName string, options ...utils.OptionFunc)
}

type Option struct {
	//jobName string
	QueueName string
	Data interface{}
}

type Job interface {
	Handle(data interface{})
}

type Payload struct {
	JobName  string
	Timeout  time.Time
	Delay    time.Time
	MaxTries int8
	Data interface{}
}

func WithQueueName(queueName string) utils.OptionFunc {
	return func(option interface{}) {
		option.(*Option).queueName = queueName
	}
}
//func WithJobName(jobName string) utils.OptionFunc {
//	return func(option interface{}) {
//		option.(*Option).jobName = jobName
//	}
//}
func WithData(data interface{}) utils.OptionFunc {
	return func(option interface{}) {
		option.(*Option).data = data
	}
}

type Manager struct {
	jobs map[string]Job
}

func NewManager() *Manager {
	return &Manager{
		jobs: make(map[string]Job),
	}
}

func (m *Manager) Register(jobName string, job Job) {
	m.jobs[jobName] = job
}

func (m *Manager) Get(jobName string) Job {
	return m.jobs[jobName]
}
