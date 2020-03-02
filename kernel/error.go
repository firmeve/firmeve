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

func (b *basicError) Render(status int, ctx contract.Context) error {
	v := map[string]interface{}{
		`status`:  status,
		`message`: b.Error(),
	}

	if ctx.Firmeve().IsDevelopment() {
		v[`meta`] = b.meta
		v[`stack`] = b.StackString()
	}

	return ctx.Render(status, v)
}

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
			stackString = append(stackString, fmt.Sprintf("%s %d %s", file, line, fn.Name()))
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
	}
}

func Errorf(format string, args ...interface{}) *basicError {
	return &basicError{
		message: fmt.Sprintf(format, args...),
		stack:   callers(),
	}
}

//@todo 这个warp重新递归包装
func ErrorWarp(err error) *basicError {
	return &basicError{
		stack:   callers(),
		err:     err,
	}

	//if e, ok := err.(*basicError); ok {
	//	return e
	//}
	//
	//return &basicError{
	//	message: err.Error(),
	//	err:     err,
	//}
}
