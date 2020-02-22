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
		entries  map[string]*contract.ContextEntity
		index    int
	}
)

func New(firmeve contract.Application, protocol contract.Protocol, handlers ...contract.ContextHandler) contract.Context {
	return &Context{
		firmeve:  firmeve,
		protocol: protocol,
		handlers: handlers,
		entries:  make(map[string]*contract.ContextEntity, 0),
		index:    0,
	}
}

func (c *Context) Firmeve() contract.Application {
	return c.firmeve
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

func (c *Context) AddEntity(key string, value interface{}) {
	c.entries[key] = &contract.ContextEntity{
		Key:   key,
		Value: value,
	}
}

func (c *Context) Entity(key string) *contract.ContextEntity {
	if v, ok := c.entries[key]; ok {
		return v
	}

	return nil
}

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
	if v, ok := key.(string); ok {
		values := c.protocol.Values()
		if value, ok := values[v]; ok {
			return value
		}
	}

	return nil
}
