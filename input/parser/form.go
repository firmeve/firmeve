package parser

import (
	"github.com/go-playground/form/v4"
	"net/url"
)

type (
	FormData = url.Values

	Form struct {
		original FormData
	}
)

var (
	formDecoder = form.NewDecoder()
)

func (j *Form) Bind(v interface{}) error {
	return formDecoder.Decode(v, j.original)
}

func (j *Form) Has(key string) bool {
	_, ok := j.original[key]
	return ok
}

func (j *Form) Get(key string) interface{} {
	return j.original.Get(key)
}

func NewForm(data FormData) *Form {
	j := &Form{
		original: data,
	}
	return j
}
