package http

//import (
//	"fmt"
//	"testing"
//)

//import (
//	"context"
//	"github.com/firmeve/firmeve/server"
//	"testing"
//	net_http "net/http"
//)
//
//func TestServer(t *testing.T) {
//	router := NewRouter()
//	router.Get(`/`, func(ctx context.Context) {
//		ctx.(*server.Context).Protocol.(*Context).String("abcdef")
//	})
//
//	net_http.ListenAndServe(":28818", router.router)
//}

//func TestRouter_Get(t *testing.T) {
//	router1 := NewRouter()
//	router1.Get(`/abc`, func(ctx *Context) {
//
//	}).Name(`abcd`)
//
//	router1.Group(func(router *Router) {
//		router.Get("/sub/ss", func(ctx *Context) {
//
//		})
//	})
//
//	router1.PathPrefix("sb").Group(func(router *Router) {
//		router.Get("/ssss", func(ctx *Context) {
//
//		}).Name("ssss")
//	})
//	//z := mux.NewRouter()
//	//z.HandleFunc(`/abc`, func(writer net_http.ResponseWriter, request *net_http.Request) {
//	//
//	//}).Name("def")
//
//	fmt.Println(router1.router.Get(`abcd`).GetMethods())
//	//fmt.Println(z.Get(`def`).GetName())
//}