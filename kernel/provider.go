package kernel

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support"
)

type (
	BaseProvider struct {
		Config      contract.Configuration `inject:"config"`
		Application contract.Application   `inject:"application"`
	}
)

func (p *BaseProvider) Bind(name string, prototype interface{}, options ...support.Option) {
	p.Application.Bind(name, prototype, options...)
}

func (p *BaseProvider) Resolve(abstract interface{}, params ...interface{}) interface{} {
	return p.Application.Make(abstract, params...)
}

func (p *BaseProvider) BindConfig(key string, object interface{}) error {
	// Application.Make(`config`).(*Config)
	return p.Config.Bind(key, object)
}
