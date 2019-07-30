package http

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}
//
//func (c *Context) Deadline() (deadline time.Time, ok bool) {
//	panic("implement me")
//}
//
//func (c *Context) Done() <-chan struct{} {
//	return nil
//}
//
//func (c *Context) Err() error {
//	return nil
//}
//
//func (c *Context) Value(key interface{}) interface{} {
//	panic("implement me")
//}
