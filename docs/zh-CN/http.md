## 简介
基础的`Http`服务和`Router`以及`Context`处理

### 路由定义
```go
//基础示例
router := http.New(firmeve.New())
router.GET("/ping", func(ctx *http.Context) {
    ctx.Write([]byte("pong"))
    ctx.Next()
})
```

### 路由中间件
```go
router.GET("/ping", func(ctx *http.Context) {
    ctx.Write([]byte("pong"))
    ctx.Next()
}).Before(func(ctx *http.Context) {
  ctx.Write([]byte("Before"))
  ctx.Next()
}).After(func(ctx *http.Context) {
   ctx.Write([]byte("After"))
   ctx.Next()
})
```

### 路由分组
```go
v1 := router.Group("/api/v1").Before(func(ctx *http.Context) {
                               ctx.Write([]byte("Before"))
                               ctx.Next()
                             }).After(func(ctx *http.Context) {
                                ctx.Write([]byte("After"))
                                ctx.Next()
                              })
{
	v1.Get("/ping", func(ctx *http.Context) {
       ctx.JSON(map[string]string{
       	    "message": "something"
       })
       ctx.Next()
   })
}
```


### 启动Http服务

#### Http服务
```bash
go run main.go http:serve --host=0.0.0.0:22182
```

#### Https服务
```bash
go run main.go http:serve --host=0.0.0.0:22182 --key-file=server.key --cert-file=server.crt
```

#### Http2服务
```bash
go run main.go http:serve --host=0.0.0.0:22182 --key-file=server.key --cert-file=server.crt --http2
```

> 假设`server.key`,`server.crt`和`main`在同一目录


### Context

todo

