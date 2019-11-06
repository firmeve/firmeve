package event

import (
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
	"testing"
)

func listenAOne(params ...interface{}) interface{} {
	return params[0].(int) + params[1].(int)
}

func listenATwo(params ...interface{}) interface{} {
	return params[0].(int) + params[1].(int)
}

func TestBaseDispatcher(t *testing.T) {
	dispatch := New()
	dispatch.Listen("a", listenAOne)
	dispatch.Listen("a", listenATwo)

	results := dispatch.Dispatch("a", 1, 2)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 3, results[0].(int))
	assert.Equal(t, 3, results[0].(int))
}

func listenBOne(params ...interface{}) interface{} {
	return false
}

func listenBTwo(params ...interface{}) interface{} {
	return params[0].(int) + params[1].(int)
}

func TestFalseDispatcher(t *testing.T) {
	dispatch := New()
	dispatch.Listen("b", listenBOne)
	dispatch.Listen("b", listenBTwo)

	results := dispatch.Dispatch("b", 1, 2)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, false, results[0].(bool))
}

func TestNamePanic(t *testing.T) {
	assert.Panics(t, func() {
		dispatch := New()
		dispatch.Dispatch("c", listenBOne)
	}, `the event not exists`)
}

func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.Instance()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("event"))
	assert.Equal(t, true, firmeve.Has(`event`))
}
