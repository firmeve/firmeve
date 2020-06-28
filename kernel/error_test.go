package kernel

import (
	"errors"
	"testing"
)

func TestErrorf(t *testing.T) {
	err := Errorf("The error %s", "something")
	ErrorWarp(errors.New("something"))

	err.SetMeta(`a`, `1`)
}
