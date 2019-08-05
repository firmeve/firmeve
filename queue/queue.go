package queue

import "time"

type Queue interface {
	Push(job Job, queue string, data interface{})
	Pop(queue string)
	Later(delay time.Time, job Job, queue string, data interface{})
}

type Job interface {
	Handle(data interface{})
}