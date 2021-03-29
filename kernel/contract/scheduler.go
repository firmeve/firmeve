package contract

type (
	Scheduler interface {
		Send(message *SchedulerMessage) error
		Run() error
		Close() error
		Decrement(workerNum int32)
		Increment(workerNum int32) error
		RegisterHandler(name string, handler SchedulerHandler)
	}

	SchedulerHandler interface {
		Handle(message *SchedulerMessage) error
	}

	SchedulerFailed interface {
		Failed(err error)
	}

	SchedulerMessage struct {
		Worker  int32
		Handler string
		Message interface{}
	}
)
