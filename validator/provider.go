package validator

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type (
	Provider struct {
		kernel.BaseProvider
	}
)

func (p Provider) Name() string {
	return `validator`
}

func (p Provider) Register() {
	validate := newValidator()
	trans := newTranslator(validate, p.Config.GetString(`framework.lang`))
	newValidator := New(validate, trans)

	p.Bind(`validator.trans`, trans)
	p.Bind(`validator`, newValidator)
}

func (p Provider) Boot() {

}

func newValidator() *validator.Validate {
	return validator.New()
}

func newTranslator(validate *validator.Validate, lang string) ut.Translator {
	var (
		translator locales.Translator
		langString string
	)
	if lang == `zh-CN` {
		translator = zh.New()
		langString = `zh`
	} else {
		translator = en.New()
		langString = `en`
	}

	trans, _ := ut.New(translator).GetTranslator(langString)

	// register default translation
	if lang == `zh-CN` {
		zh_translations.RegisterDefaultTranslations(validate, trans)
	} else {
		en_translations.RegisterDefaultTranslations(validate, trans)
	}

	return trans
}
