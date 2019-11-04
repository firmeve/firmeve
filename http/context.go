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
	index    int
}

func newContext(writer http.ResponseWriter, r *http.Request, handlers ...HandlerFunc) *Context {
	return &Context{
		request:  r,
		response: writer,
		handlers: handlers,
		index:    0,
	}
}

func (ctx *Context) Query(key string) interface{} {
	return ctx.request.URL.Query().Get(key)
}

func (ctx *Context) Form(key string) string {
	return ctx.request.FormValue(key)
}

func (ctx *Context) StatusCode(code int) *Context {
	ctx.response.WriteHeader(500)
	return ctx
}

func (ctx *Context) SetHeader(key, value string) *Context {
	ctx.response.Header().Set(key, value)
	return ctx
}

func (ctx *Context) Post(key string) string {
	return ctx.request.Form.Get(key)
}

func (ctx *Context) Write(bytes []byte) *Context {
	ctx.response.WriteHeader(http.StatusOK)
	ctx.response.Write(bytes)
	return ctx
}

func (ctx *Context) String(content string) *Context {
	ctx.response.WriteHeader(http.StatusOK)
	ctx.response.Write([]byte(content))
	return ctx
}

func (ctx *Context) Flush() *Context {
	ctx.response.(http.Flusher).Flush()
	return ctx
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
