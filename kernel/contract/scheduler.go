package contract

import "context"

type (
	SchedulerHandler interface {
		Handle(ctx context.Context, data interface{})

		ParentCtx() context.Context
	}

	Scheduler interface {
		RegisterHandler(name string, handler SchedulerHandler)

		Dispatch()

		Delivery(handler string, data interface{})
	}
)
