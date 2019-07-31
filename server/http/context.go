package http

import (
	"github.com/go-chi/chi"
	net_http "net/http"
	"time"
)

type Context struct {
	request  *net_http.Request
	response net_http.ResponseWriter
}

func (c *Context) Param(key string) string {
	return chi.URLParam(c.request, key)
}

func (c *Context) String(content string) int {
	result,err :=  c.response.Write([]byte(content))
	if err != nil{
		panic(err)
	}
	return result
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	panic("implement me")
}
