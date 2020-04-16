package validator

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
)

type validateError struct {
	message string
	errors  map[string]string
}

func Error(message string, errors map[string]string) *validateError {
	return &validateError{
		message: message,
		errors:  errors,
	}
}

func (v *validateError) Errors() map[string]string {
	return v.errors
}

func (v *validateError) Error() string {
	var firstMessage string
	for key, msg := range v.errors {
		firstMessage = fmt.Sprintf("%s:%s", key, msg)
		break
	}
	return fmt.Sprintf("%s: %s", v.message, firstMessage)
}

func (v *validateError) Render(status int, ctx contract.Context) error {
	err := map[string]interface{}{
		`status`:  status,
		`message`: v.message,
		`errors`:  v.errors,
	}

	return ctx.Render(status, err)
}
