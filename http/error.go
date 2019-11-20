package http

import (
	strings2 "github.com/firmeve/firmeve/support/strings"
	http2 "net/http"
)

type ErrorResponse interface {
	Response(writer http2.ResponseWriter)
}

type Error struct {
	code    int
	message string
	err     error
}

func (h *Error) Error() string {
	if h.err != nil {
		return strings2.Join(``, h.message, h.err.Error())
	}

	return h.message
}

func (h *Error) Response(writer http2.ResponseWriter) {
	http2.Error(writer, h.message, h.code)
}

func NewErrorWithError(code int, message string, err error) *Error {
	return &Error{
		code:    code,
		message: message,
		err:     err,
	}
}

func NewError(code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
		err:     nil,
	}
}

func Error403(message string) *Error {
	return NewError(403, message)
}

func Error500(message string) *Error {
	return NewError(500, message)
}

func Error422(message string) *Error {
	return NewError(422, message)
}
