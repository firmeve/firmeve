package http

import (
	"context"
	"time"
)

type Context struct {
	context context.Context
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	panic("implement me")
}
