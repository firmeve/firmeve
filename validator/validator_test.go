package validator

import (
	//ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	ut "github.com/go-playground/universal-translator"
	//"github.com/go-playground/validator/v10"
	"testing"
)

func app() contract.Application {
	firmeve := testing2.TestingModeFirmeve()
	firmeve.Register(new(Provider), true)
	firmeve.Boot()
	return firmeve
}

func TestValidator_RegisterValidation(t *testing.T) {
	app := app()

	v := app.Get(`validator`).(*Validator)
	err := v.RegisterTranslationValidation(`mobile`, func(fl validator.FieldLevel) bool {
		return true
	}, func(ut ut.Translator) error {
		return ut.Add(`mobile`, `手机号格式不正确`, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(`mobile`, fe.Field())
		return t
	})
	assert.Nil(t, err)
}
