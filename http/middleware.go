package http

import (
	"context"
	"errors"
	"github.com/firmeve/firmeve/logger"
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

	c.Value("context").(*Context).Firmeve.Get(`logger`).(logging.Loggable).Error(message, map[string]interface{}{
		`error`:   err,
		`context`: c.Value("context"),
	})
}

func render(err interface{}, c *Context) {
	if v, ok := err.(error); ok {
		var HttpErrorResponse ErrorResponse
		var HttpError *Error
		if errors.As(v, &HttpErrorResponse) {
			HttpErrorResponse.Response(c)
		} else if errors.As(v, &HttpError) {
			HttpError.Response(c)
		} else {
			//@todo 增加debug 判断 strings2.Join(` `, `Server error`, v.Error())
			message := `Server error`
			c.AbortWithError(http.StatusInternalServerError, message, v)
		}
	} else if v, ok := err.(string); ok {
		c.Abort(http.StatusInternalServerError, string(v))
	} else {
		c.Abort(http.StatusInternalServerError, `Server error`)
	}
}
