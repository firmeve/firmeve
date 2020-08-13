package kernel

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"testing"
)

type Logger struct {
}

func (l Logger) Debug(context ...interface{}) {
	log.Print(context...)
}

func (l Logger) Info(context ...interface{}) {
	log.Print(context...)
}

func (l Logger) Warn(context ...interface{}) {
	panic("implement me")
}

func (l Logger) Error(context ...interface{}) {
	log.Print(context...)
}

func (l Logger) Fatal(context ...interface{}) {
	panic("implement me")
}

func (l Logger) Panic(context ...interface{}) {
	panic("implement me")
}

func (l Logger) Writer(channel string) io.Writer {
	panic("implement me")
}

func (l Logger) With(context ...interface{}) contract.Loggable {
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
