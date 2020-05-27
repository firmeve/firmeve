package kernel

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support"
)

type (
	BaseProvider struct {
		Firmeve     contract.Application `inject:"application"`
		Application contract.Application `inject:"application"`
	}
)

func (p *BaseProvider) Bind(name string, prototype interface{}, options ...support.Option) {
	p.Firmeve.Bind(name, prototype, options...)
}

func (p *BaseProvider) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return p.Firmeve.Make(abstract, params...)
}
