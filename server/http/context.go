package http

import (
	net_http "net/http"
	"time"
)

type Context struct {
	request  *net_http.Request
	response net_http.ResponseWriter
}

func NewContext(response net_http.ResponseWriter, request *net_http.Request) *Context {
	return &Context{
		request:  request,
		response: response,
	}
}

//func (c *Context) Param(key string) string {
//	return c.request.P(c.request, key)
//}

func (c *Context) String(content string) int {
	result, err := c.response.Write([]byte(content))
	if err != nil {
		panic(err)
	}
	return result
}

func (c *Context) Request() *net_http.Request {
	return c.request
}

// -------------- Context interface ---------------------

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
	return nil
}
