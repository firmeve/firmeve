package contract

type (
	Scheduler interface {
		Send(message *SchedulerMessage) error
		Run() error
		Close() error
		Decrement(workerNum int32)
		Increment(workerNum int32)
		SetAvailableWorker(index int32)
		RegisterHandler(name string, handler SchedulerHandler)
		Handler(name string) SchedulerHandler
	}

	SchedulerWorker interface {
		Input() chan<- *SchedulerMessage
		Close() error
	}

	SchedulerHandler interface {
		Handle(message *SchedulerMessage) error
	}

	SchedulerFailed interface {
		Failed(err error)
	}

	SchedulerMessage struct {
		Handler string
		Message interface{}
	}
)
