package kernel

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorf(t *testing.T) {
	err := Error(fmt.Sprintf("The error %s", "something"))
	Error(errors.New("something"))

	err.SetMeta(`a`, `1`)
}
