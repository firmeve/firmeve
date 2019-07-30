package server

import (
	"context"
	"time"
)

type Context struct {
	protocol context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		protocol: ctx,
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.protocol.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.protocol.Done()
}

func (c *Context) Err() error {
	return c.protocol.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.protocol.Value(key)
}
