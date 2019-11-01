package event

import (
	"fmt"
	"sync"
)

type listenFunc func(params ...interface{}) interface{}
type listenFuncs []listenFunc

type ListenDispatcher interface {
	Listen(name string, df listenFunc)
	Dispatch(name string, params ...interface{}) []interface{}
}

type event struct {
	listeners map[string]listenFuncs
}

var (
	mutex sync.Mutex
)

func New() ListenDispatcher {
	return &event{
		listeners: make(map[string]listenFuncs, 0),
	}
}

func (e *event) Listen(name string, df listenFunc) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := e.listeners[name]; !ok {
		e.listeners[name] = make(listenFuncs, 0)
	}

	e.listeners[name] = append(e.listeners[name], df)
}

func (e *event) Dispatch(name string, params ...interface{}) []interface{} {
	var listeners listenFuncs
	var ok bool
	if listeners, ok = e.listeners[name]; !ok {
		panic(fmt.Errorf("the event %s not exists", name))
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
