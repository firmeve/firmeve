package http

import (
	"encoding/json"
	"fmt"
	"github.com/firmeve/firmeve"
	resource2 "github.com/firmeve/firmeve/converter/resource"
	"github.com/firmeve/firmeve/converter/serializer"
	"github.com/firmeve/firmeve/support/strings"
	"github.com/go-playground/form/v4"
	"net/http"
	strings2 "strings"
	"time"
)

type HandlerFunc func(c *Context)

type Params map[string]string

var (
	formDecoder = form.NewDecoder()
)

type Context struct {
	Firmeve        *firmeve.Firmeve `inject:"firmeve"`
	request        *http.Request
	responseWriter http.ResponseWriter
	handlers       []HandlerFunc
	index          int
	params         Params
	route          *Route
}

func newContext(firmeve *firmeve.Firmeve, writer http.ResponseWriter, r *http.Request, handlers ...HandlerFunc) *Context {
	return &Context{
		Firmeve:        firmeve,
		request:        r,
		responseWriter: writer,
		handlers:       handlers,
		index:          0,
		params:         make(Params, 0),
	}
}

func (c *Context) SetParams(params Params) *Context {
	c.params = params
	return c
}

func (c *Context) SetRoute(route *Route) *Context {
	c.route = route
	return c
}

func (c *Context) FormDecode(v interface{}) interface{} {
	if err := formDecoder.Decode(v, c.request.Form); err != nil {
		panic(err)
	}

	return v
}

func (c *Context) Abort(code int, message string) {
	c.AbortWithError(code, message, nil)
}

func (c *Context) AbortWithError(code int, message string, err error) {
	NewErrorWithError(code, message, err).Response(c)
}

func (c *Context) Param(key string) string {
	value, _ := c.params[key]
	return value
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) Query(key string) interface{} {
	return c.request.URL.Query().Get(key)
}

func (c *Context) Form(key string) string {
	return c.request.FormValue(key)
}

func (c *Context) Status(code int) *Context {
	c.responseWriter.WriteHeader(code)
	return c
}

func (c *Context) Header(key string) string {
	return c.request.Header.Get(key)
}

func (c *Context) SetHeader(key, value string) *Context {
	c.responseWriter.Header().Set(key, value)
	return c
}

func (c *Context) Post(key string) string {
	return c.request.Form.Get(key)
}

func (c *Context) Write(bytes []byte) *Context {
	_, err := c.responseWriter.Write(bytes)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Context) NoContent() *Context {
	c.responseWriter.WriteHeader(204)
	return c
}

func (c *Context) Created() *Context {
	c.responseWriter.WriteHeader(201)
	return c
}

func (c *Context) String(content string) *Context {
	c.Write([]byte(content))
	return c
}

func (c *Context) IsJSON() bool {
	accept := c.Header(strings.UcFirst(`Accept`))
	accepts := strings2.Split(accept, `,`)
	for _, item := range accepts {
		if item == `application/json` || item == `application+json` {
			return true
		}
	}

	return false
}

func (c *Context) JSON(content interface{}) *Context {
	c.SetHeader(`Content-Type`, `application/json`)

	str, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Write(str)

	return c
}

func (c *Context) Data(content interface{}) *Context {
	return c.JSON(serializer.NewData(content).Resolve())
}

func (c *Context) Item(resource interface{}, option *resource2.Option) *Context {
	return c.Data(resource2.NewItem(resource, option))
}

func (c *Context) Collection(resource interface{}, option *resource2.Option) *Context {
	return c.Data(resource2.NewCollection(resource, option))
}

// JSONP serializes the given struct as JSON into the responseWriter body.
// It add padding to responseWriter body to request data from a server residing in a different domain than the client.
// It also sets the Content-Type as "application/javascript".
//func (c *Context) JSONP(code int, obj interface{}) {
//	callback := c.DefaultQuery("callback", "")
//	if callback == "" {
//		c.Render(code, render.JSON{Data: obj})
//		return
//	}
//	c.Render(code, render.JsonpJSON{Callback: callback, Data: obj})
//}

func (c *Context) Redirect(location string, code int) {
	http.Redirect(c.responseWriter, c.request, location, code)
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.responseWriter, c.request, filepath)
}

func (c *Context) FileAttachment(filepath, filename string) {
	c.responseWriter.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(c.responseWriter, c.request, filepath)
}

func (c *Context) Flush() *Context {
	c.responseWriter.(http.Flusher).Flush()
	return c
}

func (c *Context) Next() {
	if c.index < len(c.handlers) {
		c.index++
		c.handlers[c.index-1](c)
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (c *Context) Done() <-chan struct{} {
	panic("implement me")
}

func (c *Context) Err() error {
	panic("implement me")
}

func (c *Context) Value(key interface{}) interface{} {
	panic("implement me")
}

//func (c *Context) ServeHttp(w http.ResponseWriter, r *http.Request) {
//}
