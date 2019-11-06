package support

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type objectTesting struct {
	id   int
	name string
}

func withId(id int) Option {
	return func(object Object) {
		object.(*objectTesting).id = id
	}
}

func withName(name string) Option {
	return func(object Object) {
		object.(*objectTesting).name = name
	}
}

func TestApplyOption(t *testing.T) {
	var (
		id   int    = 1
		name string = "name"
	)
	// Assume external calls
	object := ApplyOption(&objectTesting{}, withId(id), withName(name)).(*objectTesting)

	assert.Equal(t, id, object.id)
	assert.Equal(t, name, object.name)
}
