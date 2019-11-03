package http

import (
	"fmt"
)

type abcd func(name string) string

func (a *abcd) c() {
	fmt.Printf("%#v",a)
}
func (a abcd) d() string {
	fmt.Printf("%#v",a)
	return "a"
}

type HandlerFunc2 func(a string,r string)

// ServeHTTP calls f(w, r).
func (f HandlerFunc2) ServeHTTP(a string,r string) {
	f(a, r)
}
//
//func TestRouter_Use(t *testing.T) {
//
//	//abc := abc("a")
//	//abc.c()
//
//	//abc("abc").d()
//	//netHttp.HandlerFunc("a","b").ServeHTTP()
//	//HandlerFunc2("a","b").ServeHTTP("a","b")
//
//	//println(abcd("a"))
//	//z.c()
//
//
//	//router := New()
//	//router.GET("abc", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	//
//	//}).Use()
//
//	router := New()
//	router.GET("/abc", func(ctx *Context) {
//		ctx.Write([]byte("abc"))
//		ctx.Next()
//	}).After(func(ctx *Context) {
//		ctx.Write([]byte("SSS"))
//		ctx.Next()
//	}).After(func(ctx *Context) {
//		ctx.Write([]byte("SSSAfter"))
//		ctx.Next()
//	}).Before(func(ctx *Context) {
//		ctx.Write([]byte("abcBefore"))
//		ctx.Next()
//	})
//
//	router.GET("/def", func(ctx *Context) {
//		ctx.Write([]byte("def1"))
//		ctx.Flush()
//		ctx.Next()
//	}).After(func(ctx *Context) {
//		ctx.Write([]byte("After"))
//		ctx.Next()
//	}).Before(func(ctx *Context) {
//		ctx.Write([]byte("Before"))
//		ctx.Next()
//	})
//
//	router.NotFound(func(ctx *Context) {
//		ctx.Write([]byte("NotFound"))
//	})
//
//	v1 := router.Group("/v1").After(func(ctx *Context) {
//		ctx.Write([]byte("Group v1 After"))
//		ctx.Next()
//	}).Before(func(ctx *Context) {
//		ctx.Write([]byte("Group v1 Before"))
//		ctx.Next()
//	})
//	{
//		v1.GET("/sss", func(ctx *Context) {
//			ctx.Write([]byte("bdc"))
//			ctx.Next()
//		})
//	}
//
//	http.ListenAndServe("127.0.0.1:28082",router)
//}
