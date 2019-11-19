package errors

import (
	http2 "net/http"

	"github.com/firmeve/firmeve/http"
)

type HttpError interface {
	error
	Response()
}

type Http struct {
	Context *http.Context
	Code    int
	Message string
}

func (h *Http) Error() string {
	return h.Message
}

func (h *Http) Response() {
	http2.Error(h.Context.ResponseWriter(), h.Message, h.Code)
}

func HttpNew(code int, message string) *Http {
	return &Http{
		Code:    code,
		Message: message,
	}
}
