package memory

import (
	"fmt"
	"github.com/firmeve/firmeve/queue"
	"time"
)

type Memory struct {
	waiting chan *DataStructure
	running chan *DataStructure
}

type DataStructure struct {
	job  queue.Job
	data interface{}
}

func NewMemory() *Memory {
	return &Memory{
		waiting: make(chan *DataStructure),
		running: make(chan *DataStructure),
	}
}

func (m *Memory) Push(job queue.Job, queue string, data interface{}) {
	m.waiting <- &DataStructure{
		job:  job,
		data: data,
	}
}

func (m *Memory) Pop(queue string) *DataStructure {
	return <-m.waiting
}

func (m *Memory) Later(delay time.Time, job queue.Job, queue string, data interface{}) {
}

func (m *Memory) Run() {
	for {
		select {
		case data := <- m.waiting:
			data.job.Handle(data.data)
		case <- time.After(time.Second * 10):
			fmt.Println("超时")
		}
	}
}
