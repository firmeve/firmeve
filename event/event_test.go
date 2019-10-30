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
	dispatch := NewDispatcher()
	dispatch.listen("a", listenAOne)
	dispatch.listen("a", listenATwo)

	results := dispatch.dispatch("a", 1, 2)
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

func TestFalseDispatcher(t *testing.T)  {
	dispatch := NewDispatcher()
	dispatch.listen("b", listenBOne)
	dispatch.listen("b", listenBTwo)

	results := dispatch.dispatch("b", 1, 2)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, false, results[0].(bool))
}

func TestNamePanic(t *testing.T)  {
	assert.Panics(t, func() {
		dispatch := NewDispatcher()
		dispatch.dispatch("c", listenBOne)
	},`the event not exists`)
}


func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.NewFirmeve()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("event"))
	assert.Equal(t,true,firmeve.Has(`event`))
}


