package validator

import (
	ut "github.com/go-playground/universal-translator"
	validator2 "github.com/go-playground/validator/v10"
	"reflect"
)

type Validator struct {
	validate *validator2.Validate
	trans    ut.Translator
}

func New(validate *validator2.Validate, trans ut.Translator) *Validator {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get(`alias`)
	})

	return &Validator{
		validate: validate,
		trans:    trans,
	}
}

func (v *Validator) RegisterTranslation(tag string, registerFn validator2.RegisterTranslationsFunc, translationFn validator2.TranslationFunc) error {
	return v.validate.RegisterTranslation(tag, v.trans, registerFn, translationFn)
}

func (v *Validator) RegisterValidation(tag string, validationFunc validator2.Func) error {
	return v.validate.RegisterValidation(tag, validationFunc)
}

func (v *Validator) RegisterTranslationValidation(tag string, validationFunc validator2.Func, registerFn validator2.RegisterTranslationsFunc, translationFn validator2.TranslationFunc) error {
	if err := v.RegisterValidation(tag, validationFunc); err != nil {
		return err
	}

	return v.RegisterTranslation(tag, registerFn, translationFn)
}

func (v *Validator) Validate(val interface{}) error {
	err := v.validate.Struct(val)
	var errors = make([]*MessageStruct, 0)
	if e, ok := err.(validator2.ValidationErrors); ok {
		for _, ve := range e {
			errors = append(errors, &MessageStruct{
				Key:       ve.StructField(),
				Message:   ve.Translate(v.trans),
				Namespace: ve.StructNamespace(),
			})
		}
	}

	return Error("validation error", errors)
}
