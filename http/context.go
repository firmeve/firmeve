package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/firmeve/firmeve"
)

type HandlerFunc func(ctx *Context)

type Params map[string]string

type Context struct {
	Firmeve  *firmeve.Firmeve `inject:"firmeve"`
	request  *http.Request
	response http.ResponseWriter
	handlers []HandlerFunc
	index    int
	params   Params
	route    *Route
}

func newContext(writer http.ResponseWriter, r *http.Request, handlers ...HandlerFunc) *Context {
	firmeve := *firmeve.Instance()
	return &Context{
		Firmeve:  &firmeve,
		request:  r,
		response: writer,
		handlers: handlers,
		index:    0,
		params:   make(Params, 0),
	}
}

func (ctx *Context) SetParams(params Params) *Context {
	ctx.params = params
	return ctx
}

func (ctx *Context) SetRoute(route *Route) *Context {
	ctx.route = route
	return ctx
}

func (ctx *Context) Param(key string) string {
	value, _ := ctx.params[key]
	return value
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}

func (ctx *Context) Response(key string) http.ResponseWriter {
	return ctx.response
}

func (ctx *Context) Query(key string) interface{} {
	return ctx.request.URL.Query().Get(key)
}

func (ctx *Context) Form(key string) string {
	return ctx.request.FormValue(key)
}

func (ctx *Context) Status(code int) *Context {
	ctx.response.WriteHeader(code)
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
	ctx.response.Write(bytes)
	return ctx
}

func (ctx *Context) String(content string) *Context {
	ctx.Write([]byte(content))
	return ctx
}

func (ctx *Context) Json(content interface{}) *Context {
	ctx.SetHeader(`Content-Type`, `application/json`)
	str, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	ctx.Write(str)
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
