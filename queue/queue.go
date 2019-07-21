package queue

import "time"

type Queue interface {
	Push(job interface{}, queue string, data interface{})
	Pop(queue string)
	Later(delay time.Time, job interface{}, queue string, data interface{})
}
