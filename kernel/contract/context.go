package contract

import (
	"context"
)

type (
	ContextHandler func(c Context)

	ContextEntity struct {
		Key   string
		Value interface{}
	}

	Context interface {
		context.Context

		Protocol() Protocol

		Next()

		Handlers() []ContextHandler
	}
)


//AddEntity(entity *ContextEntity)
//Entity(key string)
//
//
//func (c *Context) Entity(key string) *entity {
//	if v, ok := c.entities[key]; ok {
//		return v
//	}
//
//	return nil
//}
//
//func (c *Context) EntityValue(key string) interface{} {
//	if v, ok := c.entities[key]; ok {
//		return v.Value
//	}
//
//	return nil
//}
//
////Output()