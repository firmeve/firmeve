package contract

import "net/http"

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
