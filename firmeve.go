package firmeve

import (
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
)

type Container interface {
	Get(id string) interface{}
	Has(id string) bool
}

type binding struct {
	shared bool
	object interface{}
}

type Firmeve struct {
	bindings map[string]binding
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[string]binding),
		}
	})

	return firmeve
}

func (f *Firmeve) bind(id string, value interface{}) {
	f.bindings[id] = newBinding(false, value)
}

// ---------------------------- binding ------------------------

func newBinding(shared bool, object interface{}) binding {
	return binding{
		shared: shared,
		object: object,
	}
}
