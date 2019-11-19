package http

import (
	http2 "net/http"
)

type IError interface {
	error
	Response()
}

type Error struct {
	Code    int
	Message string
	Context *Context
}

func (h *Error) Error() string {
	return h.Message
}

func (h *Error) Response() {
	http2.Error(h.Context.ResponseWriter(), h.Message, h.Code)
}

func NewError(code int, message string, ctx *Context) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Context: ctx,
	}
}
