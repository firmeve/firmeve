package scheduler

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
	"sync/atomic"
)

type (
	worker struct {
		s          contract.Scheduler
		message    schedulerMessageChan
		index      int32
		stopped    int32
		wait       sync.WaitGroup
		concurrent bool
	}
)

func newWorker(s contract.Scheduler, index int32, concurrent bool) contract.SchedulerWorker {
	return (&worker{
		index:      index,
		s:          s,
		stopped:    0,
		concurrent: concurrent,
		message:    make(schedulerMessageChan, 0),
	}).init()
}

func (w *worker) init() *worker {
	w.wait.Add(1)
	go func() {
		defer w.wait.Done()
		for v := range w.message {
			// goroutine hosting
			if w.concurrent {
				go w.handle(v)
			} else {
				w.handle(v)
			}
			// Queue for the next round of processing
			w.s.SetAvailableWorker(w.index)
		}
	}()

	return w
}

func (w *worker) handle(message *contract.SchedulerMessage) {
	// Pre-operation
	handler := w.s.Handler(message.Handler)
	// If not found, discard
	if handler == nil {
		return
	}
	message.Worker = w.index

	// check recover, error wrap
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

	// Usually there will be a third-party package that will return an error,
	// which only needs to be thrown in the handle, and there will be failed unified processing
	if err := handler.Handle(message); err != nil {
		if f, ok := handler.(contract.SchedulerFailed); ok {
			f.Failed(err)
		}
	}
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
