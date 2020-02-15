package event

import (
	"errors"
	"github.com/firmeve/firmeve/kernel/contract"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	result bool
	err    error
}

func (m *mockHandler) Handle(params map[string]interface{}) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}

	return params[`p1`].(int) + params[`p2`].(int), nil
}

func TestBaseDispatcher(t *testing.T) {
	dispatch := New()
	dispatch.ListenMany("a", []contract.EventHandler{
		&mockHandler{},
		&mockHandler{},
	})
	results := dispatch.Dispatch("a", map[string]interface{}{`p1`: 1, `p2`: 2})
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 3, results[0].(int))
	assert.Equal(t, 3, results[0].(int))

	nothingResult := dispatch.Dispatch("nothing", map[string]interface{}{`p1`: 1, `p2`: 2})
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

	results := dispatch.Dispatch("b", map[string]interface{}{`p1`: 1, `p2`: 2})
	assert.Equal(t, 1, len(results))
}

//func TestNamePanic(t *testing.T) {
//	assert.Panics(t, func() {
//		dispatch := New()
//		dispatch.Dispatch("c", listenBOne)
//	}, `the event not exists`)
//}
//
//func TestProvider_Register(t *testing.T) {
//	firmeve := testing2.TestingModeFirmeve()
//	//firmeve.Register(new(Provider),true)
//	//assert.Equal(t, true, firmeve.HasProvider("event"))
//	//assert.Equal(t, true, firmeve.Has(`event`))
//	//
//	//provider := firmeve.Resolve(new(Provider)).(*Provider)
//	//assert.Equal(t, firmeve, provider.Firmeve)
//}
