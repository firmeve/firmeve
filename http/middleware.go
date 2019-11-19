package http

import (
	"errors"
	"net/http"

	logging "github.com/firmeve/firmeve/logger"

	errors2 "github.com/firmeve/firmeve/http/errors"
)

var (
	httpError errors2.HttpError
)

func ServerError(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Firmeve.Get(`logger`).(logging.Loggable).Debug(err.(error).Error(), ctx)
			if errors.Is(err.(error), httpError) {
				err.(errors2.HttpError).Response()
			} else {
				ctx.Status(http.StatusInternalServerError).String("Server error")
			}
		}
	}()
	ctx.Next()
}
