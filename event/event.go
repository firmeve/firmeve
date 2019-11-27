package event

import (
	"sync"
)

type (
	InParams  map[string]interface{}
	outParams []interface{}

	Handler interface {
		Handle(params InParams) (interface{}, error)
	}

	handlers []Handler

	IDispatcher interface {
		Listen(name string, handler Handler)
		ListenMany(name string, handlers handlers)
		Dispatch(name string, params InParams) outParams
		Has(name string) bool
	}

	event struct {
		listeners map[string]handlers
	}
)

var (
	mutex sync.Mutex
)

func New() IDispatcher {
	return &event{
		listeners: make(map[string]handlers, 0),
	}
}

func (e *event) Listen(name string, handler Handler) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := e.listeners[name]; !ok {
		e.listeners[name] = make(handlers, 0)
	}

	e.listeners[name] = append(e.listeners[name], handler)
}

func (e *event) ListenMany(name string, handlerMany handlers) {
	for _, handler := range handlerMany {
		e.Listen(name, handler)
	}
}

func (e *event) Dispatch(name string, params InParams) outParams {
	if !e.Has(name) {
		return nil
	}

	results := make(outParams, 0)
	for _, listener := range e.listeners[name] {
		result, err := listener.Handle(params)
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
