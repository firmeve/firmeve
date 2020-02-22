package contract

import (
	"context"
	//"encoding/json"
	//"io"
)

type (
	ContextHandler func(c Context)

	contextEntity struct {
		Key   string
		Value interface{}
	}

	Context interface {
		context.Context

		Firmeve() Application

		Protocol() Protocol

		Next()

		Handlers() []ContextHandler

		AddEntity(key string, value interface{})

		Entity(key string) *contextEntity

		//Values() map[string][]string

		Bind(v interface{}) error

		BindWith(b Binding, v interface{}) error

		Render(v interface{}) error

		RenderWith(r Render, v interface{}) error
	}
)
