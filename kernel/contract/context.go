package contract

import (
	"context"
	//"encoding/json"
	//"io"
)

type (
	ContextHandler func(c Context)

	ContextEntity struct {
		Key   string
		Value interface{}
	}

	ContextValues map[string][]string

	Context interface {
		context.Context

		Protocol() Protocol

		Next()

		Handlers() []ContextHandler

		Bind(v interface{}) error

		BindWith(b Binding, v interface{}) error

		Render(v interface{}) error

		RenderWith(r Render, v interface{}) error
	}
)
