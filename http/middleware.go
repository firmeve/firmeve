package http

import (
	"errors"
	"net/http"

	logging "github.com/firmeve/firmeve/logger"
)

var (
	httpError IError
)

func ServerError(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Firmeve.Get(`logger`).(logging.Loggable).Debug(err.(error).Error(), ctx)
			if errors.Is(err.(error), httpError) {
				err.(IError).Response()
			} else {
				ctx.Status(http.StatusInternalServerError).String("Server error")
			}
		}
	}()
	ctx.Next()
}
