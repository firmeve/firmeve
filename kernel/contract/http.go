package contract

import "net/http"

type (
	HttpProtocol interface {
		Protocol

		Request() *http.Request

		ResponseWriter() http.ResponseWriter

		Values() map[string][]string
	}
)
