package http

import (
	"errors"
	jwt2 "github.com/firmeve/firmeve/jwt"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/gorilla/sessions"
	"net/http"
)

func Recovery(ctx contract.Context) {
	var logger = ctx.Resolve(`logger`).(contract.Loggable)
	defer kernel.RecoverCallback(logger, func(err interface{}, params ...interface{}) {
		render(err, ctx)
		ctx.Abort()
	}, ctx)

	ctx.Next()
}

func render(err interface{}, ctx contract.Context) {
	message := `Server error`
	if v, ok := err.(error); ok {
		ctx.Error(http.StatusInternalServerError, v)
	} else if v, ok := err.(string); ok {
		ctx.Error(http.StatusInternalServerError, kernel.Error(v))
	} else {
		ctx.Error(http.StatusInternalServerError, kernel.Error(message))
	}
}

func Session(ctx contract.Context) {
	httpProtocol := ctx.Protocol().(contract.HttpProtocol)
	httpProtocol.SetSession(
		NewSession(
			ctx.Resolve(`http.session.store`).(sessions.Store),
			httpProtocol.Request(),
			httpProtocol.ResponseWriter(),
		),
	)

	ctx.Next()
}

func Auth(ctx contract.Context) {
	httpProtocol := ctx.Protocol().(contract.HttpProtocol)
	token := httpProtocol.Header("Authorization")
	if token == `` {
		ctx.Error(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	jwt := ctx.Resolve(`jwt`).(contract.Jwt)

	// token valid
	if ok, err := jwt.Valid(token); !ok {
		if errors.Is(err, jwt2.ErrorExpired) {
			tokenData, _ := jwt.Refresh(token)
			token = tokenData.Token
		} else {
			ctx.Error(http.StatusUnauthorized, err)
			return
		}
	}

	// token parse
	parse, _ := jwt.Parse(token)
	ctx.AddEntity(`uid`, parse.Audience)

	ctx.Next()

	// response
	httpProtocol.ResponseWriter().Header().Set("Authorization", token)
}
