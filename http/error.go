package http

import (
	http2 "net/http"
)

type Error interface {
	error
	Response()
}

type Http struct {
	Code    int
	Message string
	Context *Context
}

func (h *Http) Error() string {
	return h.Message
}

func (h *Http) Response() {
	http2.Error(h.Context.ResponseWriter(), h.Message, h.Code)
}

func NewError(code int, message string, ctx *Context) *Http {
	return &Http{
		Code:    code,
		Message: message,
		Context: ctx,
	}
}
