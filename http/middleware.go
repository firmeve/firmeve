package http

import "net/http"

func ServerError(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Status(http.StatusInternalServerError).String("Server error")
		}
	}()
	ctx.Next()
}
