package http

import "net/http"

func ServerError(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.StatusCode(http.StatusInternalServerError).Write([]byte("Server error"))
		}
	}()
	ctx.Next()
}
