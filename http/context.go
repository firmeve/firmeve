package http

import (
	"net/http"
	"time"
)

type HandlerFunc func(ctx *Context)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	handlers []HandlerFunc
	index int
}

func newContext(writer http.ResponseWriter, r *http.Request, handlers ...HandlerFunc) *Context {
	return &Context{
		request:  r,
		response: writer,
		handlers: handlers,
		index: 0,
	}
}

func (ctx *Context) Query(key string) interface{} {
	return ctx.request.URL.Query().Get(key)
}

func (ctx *Context) Form(key string) string {
	return ctx.request.FormValue(key)
}

func (ctx *Context) Post(key string) string {
	return ctx.request.Form.Get(key)
}

func (ctx *Context) Write(bytes []byte) {
	ctx.response.Write(bytes)
}

func (ctx *Context) Flush()  {
	ctx.response.(http.Flusher).Flush()
}

func (ctx *Context) Next() {
	if ctx.index < len(ctx.handlers) {
		ctx.index++
		ctx.handlers[ctx.index-1](ctx)
	}
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

//func (ctx *Context) ServeHttp(w http.ResponseWriter, r *http.Request) {
//}