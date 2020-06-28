//go:generate mockgen -package mock -destination ../../testing/mock/mock_error.go github.com/firmeve/firmeve/kernel/contract Error,ErrorStack,ErrorRender
package contract

type Error interface {
	error

	Equal(err error) bool

	Unwrap() error

	String() string

	Meta() map[string]interface{}

	SetMeta(key string, value interface{})
}

type ErrorStack interface {
	Stack() []uintptr

	StackString() string
}

type ErrorRender interface {
	Render(status int, ctx Context) error
}
