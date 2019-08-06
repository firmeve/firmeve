package queue

type Worker struct {
	payload *Payload
}

func NewWorker(payload *Payload) *Worker {
	return &Worker{
		payload: payload,
	}
}

func (w *Worker) Handle()  {
	// 判断重试次数

	// 判断超时时间




	job := GetJob(w.payload.JobName)
	job.Handle(w.payload.Data)
}
