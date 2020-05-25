package contract

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	HttpMimeJson          = "application/json"
	HttpMimeHtml          = "text/html"
	HttpMimeXml           = "application/xml"
	HttpMimePlain         = "text/plain"
	HttpMimeForm          = "application/x-www-form-urlencoded"
	HttpMimeMultipartForm = "multipart/form-data"
	HttpMimeStream        = "application/octet-stream"
)

type (
	HttpProtocol interface {
		Protocol

		Request() *http.Request

		ResponseWriter() http.ResponseWriter

		SetHeader(key, value string)

		SetParams(params []httprouter.Param)

		Params() []httprouter.Param

		Param(key string) httprouter.Param

		SetRoute(route HttpRoute)

		Route() HttpRoute

		Header(key string) string

		SetSession(session Session)

		Session() Session

		SessionValue(key string) interface{}

		IsContentType(key string) bool

		IsAccept(key string) bool

		IsMethod(key string) bool

		ContentType() string

		Accept() []string

		SetStatus(status int)

		SetCookie(cookie *http.Cookie)

		Cookie(name string) (string, error)

		Redirect(status int, location string)
	}

	HttpRoute interface {
		Name(name string) HttpRoute

		Use(handlers ...ContextHandler) HttpRoute

		Handlers() []ContextHandler
	}

	HttpRouteGroup interface {
		Prefix(prefix string) HttpRouteGroup

		Use(handlers ...ContextHandler) HttpRouteGroup

		GET(path string, handler ContextHandler) HttpRoute

		POST(path string, handler ContextHandler) HttpRoute

		PUT(path string, handler ContextHandler) HttpRoute

		PATCH(path string, handler ContextHandler) HttpRoute

		DELETE(path string, handler ContextHandler) HttpRoute

		OPTIONS(path string, handler ContextHandler) HttpRoute

		Group(prefix string) HttpRouteGroup

		Handler(method, path string, handler http.HandlerFunc)
	}

	HttpRouter interface {
		GET(path string, handler ContextHandler) HttpRoute

		POST(path string, handler ContextHandler) HttpRoute

		PUT(path string, handler ContextHandler) HttpRoute

		PATCH(path string, handler ContextHandler) HttpRoute

		DELETE(path string, handler ContextHandler) HttpRoute

		OPTIONS(path string, handler ContextHandler) HttpRoute

		// serve static files
		Static(path string, root string) HttpRouter

		Handler(method, path string, handler http.HandlerFunc)

		HttpRouter() *httprouter.Router

		Group(prefix string) HttpRouteGroup

		ServeHTTP(w http.ResponseWriter, req *http.Request)
	}

	Session interface {
		Id() string

		Get(key string) interface{}

		GetString(key string) string
		GetInt(key string) int
		GetBool(key string) bool
		GetFloat(key string) float64

		Put(key string, value interface{}) error

		Flush()

		Delete(key string)
	}
)
