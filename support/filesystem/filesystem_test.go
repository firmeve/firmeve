package filesystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsFile(t *testing.T) {
	b := IsFile("../hash/hash.go")
	assert.Equal(t, true, b)
	b = IsFile("../hash")
	assert.Equal(t, false, b)
}
