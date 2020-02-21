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

		Values() map[string][]string

		SetHeader(key, value string)

		//IsJson() bool
		//
		//IsHtml() bool
		//
		//IsText() bool

		IsType(t string) bool
	}
)
