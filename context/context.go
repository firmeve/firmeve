package context

import (
	"fmt"
	"github.com/firmeve/firmeve/binding"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
	"sync"
	"time"
)

const (
	abortIndex = -1
)

type (
	context struct {
		application contract.Application
		protocol    contract.Protocol
		handlers    []contract.ContextHandler
		entries     map[string]*contract.ContextEntity
		index       int
	}
)

var pool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func NewContext(application contract.Application, protocol contract.Protocol, handlers ...contract.ContextHandler) contract.Context {
	ctx1, _ := application.PoolValue(`context`)
	ctx := ctx1.(*context)
	ctx.protocol = protocol
	ctx.handlers = handlers
	return ctx
}

func ReleaseContext(application contract.Application, c contract.Context) {
	application.ReleasePool(`context`, func() interface{} {
		ctx := c.(*context)
		ctx.entries = make(map[string]*contract.ContextEntity, 0)
		ctx.index = 0
		return ctx
	})
}

func (c *context) Application() contract.Application {
	return c.application
}

func (c *context) Protocol() contract.Protocol {
	return c.protocol
}

func (c *context) Error(status int, params ...interface{}) {
	paramsLen := len(params)
	var (
		err error
	)
	if paramsLen == 1 {
		if v, ok := params[0].(error); ok {
			err = v
		} else {
			err = kernel.Error(fmt.Printf("%v", params[0]))
		}
	} else { // 固定为2位参数,message,error
		err = kernel.Error(params...)
	}

	// record logger
	go c.Resolve(`logger`).(contract.Loggable).Error(
		`http error`, "error", err)

	// custom error
	if v, ok := err.(contract.ErrorRender); ok {
		if err2 := v.Render(status, c); err2 != nil {
			panic(err2)
		}
	} else {
		if err2 := render.Error.Render(c.protocol, status, err); err2 != nil {
			panic(err2)
		}
	}

	c.Abort()
}

func (c *context) Abort() {
	c.index = abortIndex
}

func (c *context) Current() int {
	return c.index
}

func (c *context) Next() {
	if c.index == abortIndex {
		return
	}

	if c.index < len(c.handlers) {
		c.index++
		c.handlers[c.index-1](c)
	}
}

func (c *context) Handlers() []contract.ContextHandler {
	return c.handlers
}

func (c *context) AddEntity(key string, value interface{}) {
	c.entries[key] = &contract.ContextEntity{
		Key:   key,
		Value: value,
	}
}

func (c *context) Entity(key string) *contract.ContextEntity {
	if v, ok := c.entries[key]; ok {
		return v
	}

	return nil
}

func (c *context) Get(key string) interface{} {
	values := c.protocol.Values()
	if value, ok := values[key]; ok {
		if len(value) == 1 {
			return value[0]
		}

		return value
	}

	return nil
}

func (c *context) Bind(v interface{}) error {
	return binding.Bind(c.protocol, v)
}

func (c *context) BindValidate(v interface{}) error {
	if err := c.Bind(v); err != nil {
		return err
	}

	return c.Resolve(`validator`).(contract.Validator).Validate(v)
}

func (c *context) BindWithValidate(b contract.Binding, v interface{}) error {
	if err := c.BindWith(b, v); err != nil {
		return err
	}

	return c.Resolve(`validator`).(contract.Validator).Validate(v)
}

func (c *context) BindWith(b contract.Binding, v interface{}) error {
	return b.Protocol(c.protocol, v)
}

func (c *context) RenderWith(status int, r contract.Render, v interface{}) error {
	return r.Render(c.protocol, status, v)
}

func (c *context) Render(status int, v interface{}) error {
	return render.Render(c.protocol, status, v)
}

func (c *context) Clone() contract.Context {
	ctxNew := new(context)
	*ctxNew = *c
	ctxNew.protocol = c.protocol.Clone()
	ctxNew.application = c.application
	ctxNew.index = c.index
	ctxNew.handlers = c.handlers

	ctxNew.entries = make(map[string]*contract.ContextEntity, len(c.entries))
	for k, v := range c.entries {
		*ctxNew.entries[k] = *v
	}
	return ctxNew
}

func (c *context) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return c.application.Make(abstract, params...)
}

// --------------------------- context.Context -> Base context ------------------------

func (c *context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *context) Done() <-chan struct{} {
	return nil
}

func (c *context) Err() error {
	return nil
}

func (c *context) Value(key interface{}) interface{} {
	if v, ok := key.(string); ok {
		return c.Get(v)
	}

	return nil
}
