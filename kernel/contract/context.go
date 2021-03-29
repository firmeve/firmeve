//go:generate mockgen -package mock -destination ../../testing/mock/mock_context.go github.com/firmeve/firmeve/kernel/contract Context
package contract

import (
	"context"
)

type (
	ContextHandler func(c Context)

	ContextEntity struct {
		Key   string
		Value interface{}
	}

	ContextPool interface {
		New(application Application) Context

		Get() Context

		Release(Context)
	}

	Context interface {
		context.Context

		Application() Application

		Protocol() Protocol

		Next()

		Handlers() []ContextHandler

		AddEntity(key string, value interface{})

		Entity(key string) *ContextEntity

		Abort()

		Current() int

		Error(status int, err ...interface{})

		Bind(v interface{}) error

		BindValidate(v interface{}) error

		BindWith(b Binding, v interface{}) error

		BindWithValidate(b Binding, v interface{}) error

		Get(key string) interface{}

		Render(status int, v interface{}) error

		RenderWith(status int, r Render, v interface{}) error

		Clone() Context

		// application make method
		Resolve(abstract interface{}, params ...interface{}) interface{}
	}
)
