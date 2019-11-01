package context

import (
	"net/http"
	"time"
)

type Context struct{
	request *http.Request
	response http.ResponseWriter
}

func (ctx *Context) Query(key string) interface{}  {
	return ctx.request.URL.Query().Get(key)
}

func (ctx *Context) Form(key string) string  {
	return ctx.request.FormValue(key)
}

func (ctx *Context) Post(key string) string  {
	return ctx.request.Form.Get(key)
}

func (ctx *Context) Next() *Context {
	return ctx
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (ctx *Context) Done() <-chan struct{} {
	panic("implement me")
}

func (ctx *Context) Err() error {
	panic("implement me")
}

func (ctx *Context) Value(key interface{}) interface{} {
	panic("implement me")
}