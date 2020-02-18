package contract

import (
	"context"
	//"encoding/json"
	//"io"
)

type (
	ContextHandler func(c Context)

	ContextEntity struct {
		Key   string
		Value interface{}
	}

	ContextValues map[string][]string

	Context interface {
		context.Context

		Protocol() Protocol

		Next()

		Handlers() []ContextHandler

		//Values() ContextValues
		Binding(v interface{})

		Render(v interface{})

	}
)

func c()  {
	//json.Marshal()
}

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