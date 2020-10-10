//go:generate mockgen -package mock -destination ../../testing/mock/mock_error.go github.com/firmeve/firmeve/kernel/contract Error,ErrorStack,ErrorRender
package contract

import "fmt"

type Error interface {
	fmt.Stringer

	error

	Equal(err error) bool

	Unwrap() error

	Meta() map[string]interface{}

	SetMeta(key string, value interface{}) Error

	Code() int

	SetCode(code int) Error

	Stack() []uintptr

	StackString() string

	Prev() error
}

type ErrorStack interface {
	Stack() []uintptr

	StackString() string
}

type ErrorRender interface {
	Render(status int, ctx Context) error
}
