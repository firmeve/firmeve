package parser

type (
	IParser interface {
		Bind(v interface{}) error
		Has(key string) bool
		Get(key string) interface{}
	}
)
