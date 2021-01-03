package scheduler

import (
	"errors"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
	"sync/atomic"
)

const defaultInitWorkersNum = 100

type (
	schedulerMessageChan chan *contract.SchedulerMessage
	workerChan           chan int32

	// The main function of the scheduler is to synchronously control N workers for concurrent processing to ensure that workers work uninterruptedly,
	// mainly to prevent excessive competition for resources by goroutines
	// Although it supports concurrent mode, it is not recommended. If you need concurrency, you can use goroutine mode.
	scheduler struct {
		workerChan      workerChan
		workerTotal     int32
		availableWorker int32
		message         schedulerMessageChan
		stopped         bool
		closeLock       sync.Mutex
		wait            sync.WaitGroup
		handlers        sync.Map
	}

	Configuration struct {
		Available int32 `json:"available" yaml:"available"`
		Worker    int32 `json:"worker" yaml:"worker"`
	}
)

func New(config *Configuration) contract.Scheduler {
	var (
		availableWorker = config.Available
		initWorkersNum  = config.Worker
	)
	if initWorkersNum == 0 {
		initWorkersNum = defaultInitWorkersNum
	}

	return (&scheduler{
		workerTotal:     initWorkersNum,
		workerChan:      make(workerChan, initWorkersNum),
		availableWorker: availableWorker,
		message:         make(schedulerMessageChan, 0),
		stopped:         false,
	}).init()
}

func (s *scheduler) RegisterHandler(name string, handler contract.SchedulerHandler) {
	s.handlers.Store(name, handler)
}

func (s *scheduler) init() *scheduler {
	// create workers
	s.createWorkers()

	return s
}

func (s *scheduler) createWorkers() {
	var i int32
	for i = 0; i < s.availableWorker; i++ {
		s.workerChan <- i
	}
}

func (s *scheduler) handle(handler contract.SchedulerHandler, message *contract.SchedulerMessage, index int32) {
	message.Worker = index

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

func (s *scheduler) Run() error {
	s.wait.Add(1)
	go func() {
		defer s.wait.Done()
		for m := range s.message {
			if handler, ok := s.handlers.Load(m.Handler); ok {
				if w, ok := <-s.workerChan; ok {
					s.wait.Add(1)
					go func(handler contract.SchedulerHandler, m *contract.SchedulerMessage, w int32) {
						defer s.wait.Done()
						defer s.setAvailableWorker(w)
						s.handle(handler, m, w)
					}(handler.(contract.SchedulerHandler), m, w)
				} else {
					break
				}
			} else {
				// Skip all handlers not found
			}
		}
	}()

	return nil
}

// Reduce workers by reducing the counter
func (s *scheduler) Decrement(workerNum int32) {
	atomic.AddInt32(&s.availableWorker, -workerNum)
	if s.availableWorker < 0 {
		atomic.StoreInt32(&s.availableWorker, 0)
	}
}

// Add worker, put index ticket into workerChan
func (s *scheduler) Increment(workerNum int32) {
	var (
		startWorker int32
		diffWorker  int32
		totalWorker int32
		i           int32
	)

	startWorker = s.availableWorker
	totalWorker = s.availableWorker + workerNum
	if totalWorker > s.workerTotal {
		diffWorker = s.workerTotal - s.availableWorker
		atomic.StoreInt32(&s.availableWorker, s.workerTotal)
	} else {
		diffWorker = workerNum
		atomic.StoreInt32(&s.availableWorker, s.availableWorker+workerNum)
	}

	for i = 0; i < diffWorker; i++ {
		s.workerChan <- (startWorker + i)
	}
}

func (s *scheduler) Send(message *contract.SchedulerMessage) error {
	if s.stopped {
		return errors.New(`scheduler is stopped`)
	}
	s.message <- message
	return nil
}

// Set each available index ticket
func (s *scheduler) setAvailableWorker(i int32) {
	if atomic.LoadInt32(&s.availableWorker) <= i {
		return
	}

	s.workerChan <- i
}

// Close
func (s *scheduler) Close() error {
	if s.stopped {
		return nil
	}

	s.closeLock.Lock()
	defer s.closeLock.Unlock()

	s.stopped = true

	// close message
	close(s.message)

	// waiting handler completed
	s.wait.Wait()

	// close
	close(s.workerChan)

	// clean
	for range s.workerChan {
	}

	return nil
}
