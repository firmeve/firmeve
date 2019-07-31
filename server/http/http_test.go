package http

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