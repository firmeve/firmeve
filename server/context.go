package server

import (
	"context"
	"time"
)

type Context struct {
	Protocol context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Protocol: ctx,
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Protocol.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.Protocol.Done()
}

func (c *Context) Err() error {
	return c.Protocol.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.Protocol.Value(key)
}
