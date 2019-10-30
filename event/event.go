package event

import (
	"fmt"
	"sync"
)

type eventFunc func(params ...interface{}) interface{}
type eventsFunc []eventFunc

type Event struct {
	listeners map[string]eventsFunc
}

var (
	event *Event
	once  sync.Once
	mu    sync.Mutex
)

func New() *Event {
	if event != nil {
		return event
	}

	once.Do(func() {
		event = &Event{
			listeners: make(map[string]eventsFunc, 0),
		}
	})

	return event
}

func (e *Event) listen(name string, df eventFunc) {
	mu.Lock()
	if _, ok := e.listeners[name]; !ok {
		e.listeners[name] = make(eventsFunc, 0)
	}

	e.listeners[name] = append(e.listeners[name], df)
	mu.Unlock()
}

func (e *Event) dispatch(name string, params ...interface{}) []interface{} {
	var listeners eventsFunc
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
