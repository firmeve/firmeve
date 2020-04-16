package validator

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	validateError struct {
		message string
		errors  []*MessageStruct
	}

	MessageStruct struct {
		Key       string `json:"key"`
		Message   string `json:"message"`
		Namespace string `json:"namespace"`
	}
)

func Error(message string, errors []*MessageStruct) *validateError {
	return &validateError{
		message: message,
		errors:  errors,
	}
}

func (v *validateError) Errors() []*MessageStruct {
	return v.errors
}

func (v *validateError) Error() string {
	var firstMessage string
	for _, msg := range v.errors {
		firstMessage = msg.Message
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
