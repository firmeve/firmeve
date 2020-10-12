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
		prev    error
		message string
		stack   []uintptr
		code    int
		meta    map[string]interface{}
	}
)

func (b *basicError) Meta() map[string]interface{} {
	return b.meta
}

func (b *basicError) SetMeta(key string, value interface{}) contract.Error {
	b.meta[key] = value
	return b
}

func (b *basicError) Error() string {
	return b.message
}

func (b *basicError) Code() int {
	return b.code
}

func (b *basicError) SetCode(code int) contract.Error {
	b.code = code
	return b
}

func (b *basicError) Equal(err error) bool {
	// TODO: 暂时使用此种方法，并非完全一致
	return errors.Is(b, err)
}

func (b *basicError) Unwrap() error {
	return b.prev
}

func (b *basicError) String() string {
	return b.Error()
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
			stackString = append(stackString, fmt.Sprintf("%s\n        %s:%d", fn.Name(), file, line))
		}
	}

	return strings.Join(stackString, "\n")
}

func (b *basicError) Prev() error {
	return b.prev
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

func Error(params ...interface{}) contract.Error {
	var (
		err     interface{}
		stacks  = callers()
		message string
		prev    error
	)
	paramsLen := len(params)
	if paramsLen == 1 {
		err = params[0]
	} else if paramsLen == 2 {
		message = params[0].(string)
		err = params[1]
	}

	// 其实和Exception一样，每个Error都会挂载上一个Error形成调用链条
	switch v := err.(type) {
	case contract.Error:
		prev = v.(error)
		stacks = append(v.Stack(), stacks...)
		if message != `` {
			message = v.(error).Error()
		}
	case error:
		prev = v
		if message != `` {
			message = v.Error()
		}
	case string:
		message = v
	default:
		if message != `` {
			message = fmt.Sprintf("%v", v)
		}
	}

	return &basicError{
		prev:    prev,
		message: message,
		stack:   stacks,
	}
}
