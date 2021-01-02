package scheduler

import (
	"errors"
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
	"sync/atomic"
)

const defaultInitWorkersNum = 300

type (
	schedulerMessageChan chan *contract.SchedulerMessage
	workerChan           chan int32

	scheduler struct {
		workers         []contract.SchedulerWorker
		workerNum       int32
		workerChan      workerChan
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
		workers:         make([]contract.SchedulerWorker, initWorkersNum),
		workerNum:       int32(initWorkersNum),
		workerChan:      make(workerChan, initWorkersNum),
		availableWorker: availableWorker,
		message:         make(schedulerMessageChan, 0),
		stopped:         false,
	}).init()
}

func (s *scheduler) RegisterHandler(name string, handler contract.SchedulerHandler) {
	s.handlers.Store(name, handler)
}

func (s *scheduler) Handler(name string) contract.SchedulerHandler {
	if v, ok := s.handlers.Load(name); ok {
		return v.(contract.SchedulerHandler)
	}
	return nil
}

func (s *scheduler) init() *scheduler {
	// create workers
	s.createWorkers()

	return s
}

func (s *scheduler) createWorkers() {
	var i int32
	for i = 0; i < s.workerNum; i++ {
		s.workers[i] = newWorker(s, i)
		if i < s.availableWorker {
			s.workerChan <- i
		}
	}
}

func (s *scheduler) Run() error {
	s.wait.Add(1)
	go func() {
		defer s.wait.Done()
		for m := range s.message {
			if w, ok := <-s.workerChan; ok {
				s.workers[w].Input() <- m
			} else {
				break
			}
		}
	}()

	return nil
}

// Reduce workers by reducing the counter
//
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
	if totalWorker > s.workerNum {
		diffWorker = s.workerNum - s.availableWorker
		atomic.StoreInt32(&s.availableWorker, s.workerNum)
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
func (s *scheduler) SetAvailableWorker(i int32) {

	if atomic.LoadInt32(&s.availableWorker) <= i {
		return
	}

	//defer func() {
	//	recover()
	//}()
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

	// stop all worker
	for i := range s.workers {
		s.workers[i].Close()
	}

	// close
	close(s.workerChan)

	return nil
}
