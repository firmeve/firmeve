package maps

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeInterface(t *testing.T) {
	m1 := map[string]interface{}{
		`a`: `1`,
		`b`: `b`,
		`c`: `c`,
	}
	m2 := map[string]interface{}{
		`c`: 10,
		`d`: true,
		`e`: `some`,
	}
	m := MergeInterface(m1, m2)
	assert.Equal(t, map[string]interface{}{
		`a`: `1`,
		`b`: `b`,
		`c`: 10,
		`d`: true,
		`e`: `some`,
	}, m)
}
