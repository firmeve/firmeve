package event

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
)

type (
	event struct {
		listeners map[string][]contract.EventHandler
	}
)

var (
	mutex sync.Mutex
)

func New() contract.Event {
	return &event{
		listeners: make(map[string][]contract.EventHandler, 0),
	}
}

func (e *event) Listen(name string, handlers ...contract.EventHandler) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := e.listeners[name]; !ok {
		e.listeners[name] = make([]contract.EventHandler, 0)
	}

	e.listeners[name] = append(e.listeners[name], handlers...)
}

func (e *event) Dispatch(name string, params ...interface{}) []interface{} {
	if !e.Has(name) {
		return nil
	}

	results := make([]interface{}, 0)
	for _, listener := range e.listeners[name] {
		result, err := listener.Handle(params...)
		if err != nil {
			break
		}
		results = append(results, result)
	}

	return results
}

func (e *event) Has(name string) bool {
	_, ok := e.listeners[name]
	return ok
}
