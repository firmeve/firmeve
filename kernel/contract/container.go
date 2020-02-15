package contract

import "github.com/firmeve/firmeve/support"

type (
	Container interface {
		Has(name string) bool
		Get(name string) interface{}
		Bind(name string, prototype interface{}, options ...support.Option)
		Make(abstract interface{}, params ...interface{}) interface{}
		Remove(name string)
		Flush()
	}
)
