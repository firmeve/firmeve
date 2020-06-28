package kernel

import (
	"errors"
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"runtime"
	"strings"
)

type (
	basicError struct {
		err     error
		message string
		stack   []uintptr
		meta    map[string]interface{}
	}
)

func (b *basicError) Meta() map[string]interface{} {
	return b.meta
}

func (b *basicError) SetMeta(key string, value interface{}) {
	b.meta[key] = value
}

func (b *basicError) Error() string {
	if b.err != nil {
		return strings.Join([]string{b.message, b.err.Error()}, `: `)
	}

	return b.message
}

func (b *basicError) Equal(err error) bool {
	return errors.Is(b, err)
}

func (b *basicError) Unwrap() error {
	return b.err
}

func (b *basicError) String() string {
	return b.Error()
}

//func (b *basicError) Render(status int, ctx contract.Context) error {
//	v := map[string]interface{}{
//		`status`:  status,
//		`message`: b.Error(),
//	}
//
//	if ctx.Firmeve().IsDevelopment() {
//		v[`meta`] = b.meta
//		v[`stack`] = b.StackString()
//	}
//	//
//	//if protocol, ok := ctx.Protocol().(contract.HttpProtocol); ok {
//	//	if protocol.IsContentType(contract.HttpMimeJson)
//	//}
//
//	return ctx.Render(status, v)
//}

func (b *basicError) Stack() []uintptr {
	return b.stack
}

func (b *basicError) StackString() string {
	var stackString []string
	for i := range b.stack {
		fn := runtime.FuncForPC(b.stack[i])
		if fn == nil {
			stackString = append(stackString, "unknown")
		} else {
			file, line := fn.FileLine(b.stack[i])
			stackString = append(stackString, fmt.Sprintf("%s\n        %s:%d", fn.Name(), file, line))
		}
	}

	return strings.Join(stackString, "\n")
}

func callers() []uintptr {
	//pc := make([]uintptr, 0)
	//n := runtime.Callers(2, pc)
	//frames := runtime.CallersFrames(pc[:n])

	//var pcs []uintptr
	////pcs := make([]uintptr, 0)
	//n := runtime.Callers(0, pcs)
	//fmt.Println(n)
	//return pcs[0:n]

	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

func Error(message string) *basicError {
	return &basicError{
		message: message,
		stack:   callers(),
		meta:    make(map[string]interface{}, 0),
	}
}

func Errorf(format string, args ...interface{}) *basicError {
	return Error(fmt.Sprintf(format, args...))
}

func ErrorWarp(err error) *basicError {
	var stacks = make([]uintptr, 0)
	if v, ok := err.(contract.ErrorStack); ok {
		stacks = append(v.Stack(), callers()...)
	} else {
		stacks = callers()
	}
	return &basicError{
		stack: stacks,
		err:   err,
		meta:  make(map[string]interface{}, 0),
	}
}
