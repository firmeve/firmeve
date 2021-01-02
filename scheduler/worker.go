package scheduler

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
	"sync/atomic"
)

type (
	worker struct {
		s       contract.Scheduler
		message schedulerMessageChan
		index   int32
		stopped int32
		wait    sync.WaitGroup
	}
)

func newWorker(s contract.Scheduler, index int32) contract.SchedulerWorker {
	return (&worker{
		index:   index,
		s:       s,
		stopped: 0,
		message: make(schedulerMessageChan, 0),
	}).init()
}

func (w *worker) init() *worker {
	w.wait.Add(1)
	go func() {
		defer w.wait.Done()
		for v := range w.message {
			w.handle(v)
		}
	}()

	return w
}

func (w *worker) handle(message *contract.SchedulerMessage) {
	handler := w.s.Handler(message.Handler)
	if handler == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			if f, ok := handler.(contract.SchedulerFailed); ok {
				var (
					e  error
					ok bool
				)
				if e, ok = err.(error); ok {
				} else {
					e = kernel.SchedulerRecoverError{
						Message: `handle panic`,
						Recover: err,
					}
				}
				f.Failed(e)
			}
		}
	}()

	if err := handler.Handle(message); err != nil {
		if f, ok := handler.(contract.SchedulerFailed); ok {
			f.Failed(err)
		}
	}

	w.s.SetAvailableWorker(w.index)
}

func (w *worker) Input() chan<- *contract.SchedulerMessage {
	return w.message
}

func (w *worker) Close() error {
	if atomic.LoadInt32(&w.stopped) == 0 {
		close(w.message)
		w.wait.Wait()
		atomic.CompareAndSwapInt32(&w.stopped, 0, 1)
	}
	return nil
}
