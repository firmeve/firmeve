package http

import (
	"errors"
	jwt2 "github.com/firmeve/firmeve/jwt"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/strings"
	"github.com/gorilla/sessions"
	"net/http"
)

func Recovery(ctx contract.Context) {
	defer panicRecovery(ctx)
	ctx.Next()
}

func panicRecovery(ctx contract.Context) {
	if err := recover(); err != nil {
		render(err, ctx)

		go report(err, ctx.Clone())
	}
}

func report(err interface{}, ctx contract.Context) {
	var (
		message string
	)
	if v, ok := err.(error); ok {
		message = v.Error()
	} else if v, ok := err.(string); ok {
		message = v
	} else {
		message = `mixed type`
	}

	//@todo 这里有问题，后续优化
	ctx.Resolve(`logger`).(contract.Loggable).Error(
		strings.Join(` `, message, "Context: %s", "Error: %#v"),
		"ctx",
		ctx,
		"error",
		err,
	)
}

func render(err interface{}, ctx contract.Context) {
	message := `Server error`
	if v, ok := err.(error); ok {
		ctx.Error(http.StatusInternalServerError, v)
	} else if v, ok := err.(string); ok {
		ctx.Error(http.StatusInternalServerError, kernel.Errorf(v))
	} else {
		ctx.Error(http.StatusInternalServerError, kernel.Errorf(message))
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
