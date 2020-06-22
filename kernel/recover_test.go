package kernel

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type Logger struct {
}

func (l Logger) Debug(message string, context ...interface{}) {
	log.Print(message)
}

func (l Logger) Info(message string, context ...interface{}) {
	log.Print(message)
}

func (l Logger) Warn(message string, context ...interface{}) {
	panic("implement me")
}

func (l Logger) Error(message string, context ...interface{}) {
	log.Print(message)
}

func (l Logger) Fatal(message string, context ...interface{}) {
	panic("implement me")
}

func (l Logger) Panic(message string, context ...interface{}) {
	panic("implement me")
}

func TestRecover(t *testing.T) {
	assert.Panics(t, func() {
		panicFunc()
	}, `error`)
}

func panicFunc() {
	defer Recover(&Logger{})
	panic("error")
}
