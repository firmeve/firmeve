package http

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/strings"
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
	ctx.Firmeve().Get(`logger`).(contract.Loggable).Error(
		strings.Join(` `, message, "Context: %s", "Error: %#v"),
		fmt.Sprintf("%#v", ctx),
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
