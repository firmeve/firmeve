package context

import (
	"github.com/firmeve/firmeve/binding"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
	"time"
)

type (
	Context struct {
		firmeve  contract.Application
		protocol contract.Protocol
		handlers []contract.ContextHandler
		index    int
	}
)

func New(firmeve contract.Application, protocol contract.Protocol, handlers ...contract.ContextHandler) contract.Context {
	return &Context{
		firmeve:  firmeve,
		protocol: protocol,
		handlers: handlers,
		index:    0,
	}
}

func (c *Context) Protocol() contract.Protocol {
	return c.protocol
}

func (c *Context) Next() {
	if c.index < len(c.handlers) {
		c.index++
		c.handlers[c.index-1](c)
	}
}

func (c *Context) Handlers() []contract.ContextHandler {
	return c.handlers
}

//
//func (c *Context) Values() map[string][]string {
//	return c.protocol.Values()
//}

func (c *Context) Bind(v interface{}) error {
	return binding.Bind(c.protocol, v)
}

func (c *Context) BindWith(b contract.Binding, v interface{}) error {
	return b.Protocol(c.protocol, v)
}

func (c *Context) RenderWith(r contract.Render, v interface{}) error {
	return r.Render(c.protocol, v)
}

func (c *Context) Render(v interface{}) error {
	return render.Render(c.protocol, v)
}

// --------------------------- context.Context -> Base context ------------------------

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	//if v, ok := key.(string); ok {
	//	return c.protocol.Values(v)
	//}

	return nil
}
