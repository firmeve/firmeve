package event

import (
	"errors"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	result bool
	err    error
}

func (m *mockHandler) Handle(app contract.Application, params ...interface{}) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}

	return params[0].(int) + params[1].(int), nil
}

func TestBaseDispatcher(t *testing.T) {
	dispatch := New(kernel.New())
	dispatch.Listen("a",
		&mockHandler{},
		&mockHandler{},
	)
	results := dispatch.Dispatch("a", 1, 2)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 3, results[0].(int))
	assert.Equal(t, 3, results[0].(int))

	nothingResult := dispatch.Dispatch("nothing", 1, 2)
	assert.Nil(t, nothingResult)
}

func TestErrorDispatcher(t *testing.T) {
	application := kernel.New()
	dispatch := New(application)
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
