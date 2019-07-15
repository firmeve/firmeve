package firmeve

import (
	"reflect"
	"strings"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
)

type Container interface {
	Get(id string) interface{}
	Has(id string) bool
	//Bind(id interface{})
}

type binding struct {
	shared bool
	object interface{}
}

type Firmeve struct {
	bindings map[string]*binding
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[string]*binding),
		}
	})

	return firmeve
}

func (f *Firmeve) Bind(id interface{}) { //, value interface{}
	//f.bindings[id] = newBinding(false, value)
	reflectType := reflect.TypeOf(id)
	if reflectType.Kind() == reflect.Ptr {
		fullName := strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, ".")
		f.bindings[fullName] = newBinding(true, id)
	}

}

func (f *Firmeve) Resolve(id interface{}) interface{}{
	reflectType := reflect.TypeOf(id)
	if reflectType.Kind() == reflect.Func {
		fullName := strings.Join([]string{reflectType.In(0).Elem().PkgPath(), reflectType.In(0).Elem().Name()}, ".")
		return reflect.ValueOf(id).Call([]reflect.Value{reflect.ValueOf(f.bindings[fullName].object)})[0].Interface()
	}

	return  nil
}

// ---------------------------- binding ------------------------

func newBinding(shared bool, object interface{}) *binding {
	return &binding{
		shared: shared,
		object: object,
	}
}
