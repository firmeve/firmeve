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

		//SetRoute(route *http2.Route)
		//
		//Route() *http2.Route

		Header(key string) string

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
)
