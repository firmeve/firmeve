package event

import (
	"fmt"
	"sync"
)

type dispatchFunc func(params ...interface{}) interface{}
type dispatchFuncs []dispatchFunc

type Dispatcher struct {
	listeners map[string]dispatchFuncs
}

var (
	event *Dispatcher
	once  sync.Once
	mu    sync.Mutex
)

func NewDispatcher() *Dispatcher {
	if event != nil {
		return event
	}

	once.Do(func() {
		event = &Dispatcher{
			listeners: make(map[string]dispatchFuncs, 0),
		}
	})

	return event
}

func (e *Dispatcher) listen(name string, df dispatchFunc) {
	mu.Lock()
	if _, ok := e.listeners[name]; !ok {
		e.listeners[name] = make(dispatchFuncs, 0)
	}

	e.listeners[name] = append(e.listeners[name], df)
	mu.Unlock()
}

func (e *Dispatcher) dispatch(name string, params ...interface{}) []interface{} {
	var listeners dispatchFuncs
	var ok bool
	if listeners, ok = e.listeners[name]; !ok {
		panic(fmt.Sprintf("the event %s not exists", name))
	}

	results := make([]interface{}, 0)
	for _, dispatchFunc := range listeners {
		result := dispatchFunc(params...)
		if v, ok := result.(bool); ok && v == false {
			results = append(results, v)
			break
		}
		results = append(results, result)
	}
	return results
}
