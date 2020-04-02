package validator

import (
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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
	config := p.Resolve(`config`).(*config2.Config).Item(`framework`)

	//fmt.Println(newValidator(), newTranslator(config.GetString(`lang`)))
	validate := newValidator()
	newValidator := New(validate, newTranslator(validate, config.GetString(`lang`)))
	//fmt.Println(newValidator)
	//fmt.Println(newValidator)
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
	}

	trans, _ := ut.New(translator).GetTranslator(langString)

	return trans
}
