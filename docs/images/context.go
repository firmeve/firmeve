package context

import (
	"context"
	"github.com/firmeve/firmeve/kernel"
	"time"
)

type (
	Context interface {
		context.Context

		Protocol() Protocol

		Next()

		Handlers() []HandlerFunc
	}

	HandlerFunc func(c *Context)

	BaseContext struct {
		firmeve  *kernel.IApplication
		protocol Protocol
		handlers []HandlerFunc
	}
)

func (c *BaseContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *BaseContext) Done() <-chan struct{} {
	return nil
}

func (c *BaseContext) Err() error {
	return nil
}

func (c *BaseContext) Value(key interface{}) interface{} {
	if v,ok := key.(string); ok {
		return c.protocol.Value(v)
	}

	return nil
}

func New()  {
	//http.Res
}
