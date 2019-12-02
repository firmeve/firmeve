package parser

import "github.com/go-playground/form/v4"

type (
	IParser interface {
		Bind(v interface{}) error
		Has(key string) bool
		Get(key string) interface{}
	}
)

var (
	FormDecoder = form.NewDecoder()
)