package path

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunFile(t *testing.T) {
	assert.Contains(t, RunFile(), "path_test")
}

func TestRunDir(t *testing.T) {
	assert.Contains(t, RunDir(), "support/path")
}

func TestRunRelative(t *testing.T) {
	assert.Contains(t, RunRelative("../"), "support")
}
