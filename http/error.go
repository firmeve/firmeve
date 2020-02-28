package http

import (
	"errors"
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"strings"
)

type (
	errorHttp struct {
		//code    int
		message string
		err     error
		status  int
		details []interface{}
	}
)

//func (h *errorHttp) Code() int {
//	return h.code
//}

func (h *errorHttp) Details() []interface{} {
	return nil
}

// Response text/plain or text/html error
func (h *errorHttp) String() string {
	return h.Error()
}

func (h *errorHttp) Equal(err error) bool {
	return errors.Is(h, err)
}

func (h *errorHttp) Status() int {
	return h.status
}

func (h *errorHttp) Response(c contract.Context) error {
	//@todo 暂时这样写
	var err = make(map[string]interface{}, 0)
	err[`message`] = h.message
	err[`status`] = h.status
	if !c.Firmeve().IsProduction() {
		err[`details`] = h.details
		err[`err`] = h.err
	}

	return c.Render(err)
}

func (h *errorHttp) Error() string {
	if h.err != nil {
		return strings.Join([]string{h.message, h.err.Error()}, ` `)
	}

	return h.message
}

func (h *errorHttp) Unwrap() error {
	return h.err
}

//func (h *Error) Response(c *Context) {
//	if c.IsJSON() {
//		c.Status(h.code).JSON(map[string]interface{}{
//			"message": h.Error(),
//		})
//	} else {
//		http2.Error(c.ResponseWriter, h.Error(), h.code)
//	}
//}

func Error(status int, message string) contract.HttpError {
	return &errorHttp{
		status:  status,
		message: message,
	}
}

func Errorf(status int, format string, args ...interface{}) contract.HttpError {
	err := fmt.Errorf(format, args)
	return &errorHttp{
		status:  status,
		message: err.Error(),
		details: args,
		err:     errors.Unwrap(err),
	}
}

//func ErrorWithError(status int, message string, err error) contract.HttpError {
//	return &errorHttp{
//		status:  status,
//		message: message,
//		err:     err,
//	}
//}

//func ErrorWithError(status int, message string, args ...interface{}) {
//	os.PathError{}
//	return &errorHttp{
//		status:  status,
//		message: message,
//		details: args,
//	}
//}
//
//func NewErrorWithError(code int, message string, err error) *Error {
//	return &Error{
//		code:    code,
//		message: message,
//		err:     err,
//	}
//}
//
//func NewError(code int, message string) *Error {
//	return &Error{
//		code:    code,
//		message: message,
//		err:     nil,
//	}
//}
//
//func Error400(message string) *Error {
//	os.PathError{}
//	http.StatusInternalServerError
//	return NewError(400, message)
//}
//
//func Error400WithError(message string, err error) *Error {
//	return NewErrorWithError(400, message, err)
//}
//
//func Error403(message string) *Error {
//	return NewError(403, message)
//}
//
//func Error403WithError(message string, err error) *Error {
//	return NewErrorWithError(403, message, err)
//}
//
//func Error422(message string) *Error {
//	return NewError(422, message)
//}
//
//func Error422WithError(message string, err error) *Error {
//	return NewErrorWithError(422, message, err)
//}
//
//func Error500(message string) *Error {
//	return NewError(500, message)
//}
//
//func Error500WithError(message string, err error) *Error {
//	return NewErrorWithError(500, message, err)
//}
