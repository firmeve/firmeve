package grpc

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"google.golang.org/grpc"
)

//type (
//	Router struct {
//		Firmeve  contract.Application
//		handlers []Server
//		services []interface{}
//	}
//
//	Server func(server *grpc.Server, service interface{})
//)
//
//func (r *Router) Register(handler Server, service interface{}) {
//	r.handlers = append(r.handlers, handler)
//	r.services = append(r.services, service)
//}
//
//func (r *Router) Handle(server *grpc.Server) {
//	for i := range r.handlers {
//		r.handlers[i](server, r.services[i])
//	}
//}

type (
	Router struct {
		Firmeve  contract.Application
		handlers []Server
	}

	Server func(server *grpc.Server)
)

func (r *Router) Register(handler Server) {
	r.handlers = append(r.handlers, handler)
}

func (r *Router) Handle(server *grpc.Server) {
	for i := range r.handlers {
		r.handlers[i](server)
	}
}

func NewRouter(firmeve contract.Application) *Router {
	return &Router{
		Firmeve:  firmeve,
		handlers: make([]Server, 0),
	}
}
