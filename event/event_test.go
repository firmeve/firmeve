package event

import (
	"errors"
	"fmt"
	"testing"

	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	result bool
	err    error
}

func (m *mockHandler) Handle(params ...interface{}) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}

	return params[0].(int) + params[1].(int), nil
}

type Func func(params ...interface{})

func (f Func) Handle(params ...interface{}) interface{} {
	fmt.Println(params...)
	return nil
}

func TestBaseDispatcher(t *testing.T) {
	dispatch := New()
	dispatch.ListenMany("a", []Handler{
		&mockHandler{},
		&mockHandler{},
	})
	results := dispatch.Dispatch("a", 1, 2)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 3, results[0].(int))
	assert.Equal(t, 3, results[0].(int))

	nothingResult := dispatch.Dispatch("nothing", 1, 2)
	assert.Nil(t, nothingResult)
}

func TestErrorDispatcher(t *testing.T) {
	dispatch := New()
	dispatch.Listen("b", &mockHandler{
		result: true,
	})
	dispatch.Listen("b", &mockHandler{
		result: false,
		err:    errors.New("error testing"),
	})
	dispatch.Listen("b", &mockHandler{
		result: true,
	})

	results := dispatch.Dispatch("b", 1, 2)
	assert.Equal(t, 1, len(results))
}

//func TestNamePanic(t *testing.T) {
//	assert.Panics(t, func() {
//		dispatch := New()
//		dispatch.Dispatch("c", listenBOne)
//	}, `the event not exists`)
//}

func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.New()
	firmeve.Register(firmeve.Make(new(Provider)).(firmeve2.Provider))
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("event"))
	assert.Equal(t, true, firmeve.Has(`event`))

	provider := firmeve.Resolve(new(Provider)).(*Provider)
	assert.Equal(t, firmeve, provider.Firmeve)
}
