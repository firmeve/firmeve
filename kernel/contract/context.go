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

	Context interface {
		context.Context

		Firmeve() Application

		Protocol() Protocol

		Next()

		Handlers() []ContextHandler

		AddEntity(key string, value interface{})

		Entity(key string) *ContextEntity

		Abort()

		//SetStatus(status int)
		//
		//Status() int

		Error(status int, err error)

		Bind(v interface{}) error

		BindWith(b Binding, v interface{}) error

		Render(status int, v interface{}) error

		RenderWith(status int, r Render, v interface{}) error
	}
)
