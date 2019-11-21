package http

import (
	http2 "net/http"
)

type ErrorResponse interface {
	Response(c *Context)
}

type Error struct {
	code    int
	message string
	err     error
}

func (h *Error) Error() string {
	return h.message
}

func (h *Error) Unwrap() error {
	return h.err
}

func (h *Error) Response(c *Context) {
	if c.IsJSON() {
		c.Status(h.code).JSON(map[string]interface{}{
			"message": h.Error(),
		})
	} else {
		http2.Error(c.responseWriter, h.Error(), h.code)
	}
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

func Error400(message string) *Error {
	return NewError(400, message)
}

func Error400WithError(message string, err error) *Error {
	return NewErrorWithError(400, message, err)
}

func Error403(message string) *Error {
	return NewError(403, message)
}

func Error403WithError(message string, err error) *Error {
	return NewErrorWithError(403, message, err)
}

func Error422(message string) *Error {
	return NewError(422, message)
}

func Error422WithError(message string, err error) *Error {
	return NewErrorWithError(422, message, err)
}

func Error500(message string) *Error {
	return NewError(500, message)
}

func Error500WithError(message string, err error) *Error {
	return NewErrorWithError(500, message, err)
}
