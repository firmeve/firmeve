package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/firmeve/firmeve/logger"
	"github.com/firmeve/firmeve/support/strings"
	"net/http"
)

func Recovery(c *Context) {
	defer panicRecovery(c)
	c.Next()
}

func panicRecovery(c *Context) {
	if err := recover(); err != nil {
		render(err, c)

		// @todo context未测试，也可以不使用context，但 http.Context 一定只能是只读
		go report(err, context.WithValue(context.Background(), "context", c))
	}
}

func report(err interface{}, c context.Context) {
	var (
		message string
	)
	if v, ok := err.(error); ok {
		message = v.Error()
	} else if v, ok := err.(string); ok {
		message = string(v)
	} else {
		message = `mixed type`
	}

	ctx := c.Value("context").(*Context)
	ctx.Firmeve.Get(`logger`).(logging.Loggable).Error(
		strings.Join(` `, message, "Context: %s", "Error: %#v"),
		fmt.Sprintf("%#v", *ctx),
		err,
	)
}

func render(err interface{}, c *Context) {
	message := `Server error`
	if v, ok := err.(error); ok {
		var HttpErrorResponse ErrorResponse
		var HttpError *Error
		if errors.As(v, &HttpErrorResponse) {
			HttpErrorResponse.Response(c)
		} else if errors.As(v, &HttpError) {
			HttpError.Response(c)
		} else {
			//@todo 增加debug 判断 strings2.Join(` `, `Server error`, v.Error())
			NewError(http.StatusInternalServerError, message).Response(c)
		}
	} else if v, ok := err.(string); ok {
		NewError(http.StatusInternalServerError, string(v)).Response(c)
	} else {
		NewError(http.StatusInternalServerError, message).Response(c)
	}
}
