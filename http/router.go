package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type Router struct {
	router    *httprouter.Router
	routes    map[string]*Route
	routeKeys []string
}

func New() *Router {
	return &Router{
		router:    httprouter.New(),
		routes:    make(map[string]*Route, 0),
		routeKeys: make([]string, 0),
	}
}

// 只是用到httpRouter的存储路由以及查找方法，因为暂时不会前缀树算法
// 其它router,middleware这些都是我自己实现，惟一对接的就是无缝的写入一套httprouter规则的路由（后期替换为自己的路由）
// 通过ServerHttp去查找匹配路由

func (r *Router) GET(path string, handler HandlerFunc) *Route {
	return r.createRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) *Route {
	return r.createRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) *Route {
	return r.createRoute(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc) *Route {
	return r.createRoute(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) *Route {
	return r.createRoute(http.MethodDelete, path, handler)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc) *Route {
	return r.createRoute(http.MethodOptions, path, handler)
}
//
//func (r *Router) Group(method string, path string, handler HandlerFunc) *Route {
//	key := r.routeKey(method, path)
//	r.routes[key] = newRoute(path, handler)
//
//	//Only http router
//	r.router.Handler(method, path, r)
//
//	return r.routes[key]
//}

func (r *Router) createRoute(method string, path string, handler HandlerFunc) *Route {
	key := r.routeKey(method, path)
	r.routes[key] = newRoute(path, handler)

	//Only http router
	r.router.Handler(method, path, r)

	return r.routes[key]
}

func (r *Router) routeKey(method, path string) string {
	return strings.Join([]string{method, path}, `.`)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if route, ok := r.routes[r.routeKey(req.Method, req.URL.Path)]; ok {
		handlers := append(append(route.beforeHandlers, route.handler), route.afterHandlers...)
		newContext(w, req, handlers...).Next()
		return
	}

	r.router.NotFound.ServeHTTP(w, req)
}