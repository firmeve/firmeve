package scheduler

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"hash/crc32"
	"sync"
)

type (
	Scheduler struct {
		length   uint32
		queues   []chan schan
		handlers *sync.Map
	}

	Configuration struct {
		Size uint32 `json:"size" yaml:"size"`
	}

	schan struct {
		name string
		data interface{}
	}
)

func New(config *Configuration) contract.Scheduler {
	if config.Size == 0 {
		config.Size = 10
	}
	queues := make([]chan schan, config.Size)
	for i := range queues {
		queues[i] = make(chan schan)
	}
	return &Scheduler{
		length:   config.Size,
		queues:   queues,
		handlers: &sync.Map{},
	}
}

func (s *Scheduler) RegisterHandler(name string, handler contract.SchedulerHandler) {
	s.handlers.Store(name, handler)
}

func (s *Scheduler) Dispatch() {
	for i := range s.queues {
		go s.dispatchOne(s.queues[i])
	}
}

func (s *Scheduler) dispatchOne(queue <-chan schan) {
	for v := range queue {
		handler, ok := s.handlers.Load(v.name)
		// ignore not found handler
		if ok {
			currentHandler := handler.(contract.SchedulerHandler)
			// No goroutine distribution, if you need to distribute it inside the processor
			currentHandler.Handle(currentHandler.ParentCtx(), v.data)
		}
	}
}

// Delivery a chan
func (s *Scheduler) Delivery(handler string, data interface{}) {
	n := crc32.ChecksumIEEE([]byte(handler))
	s.queues[n%s.length] <- schan{
		name: handler,
		data: data,
	}
}
