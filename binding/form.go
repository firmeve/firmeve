package binding

import (
	"errors"

	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/go-playground/form/v4"
)

type (
	Form struct {
	}
)

var (
	ProtocolTypeError = errors.New("the protocol only support http protocol")
	formDecoder       = form.NewDecoder()
)

func (f *Form) Name() string {
	return `form`
}

func (f *Form) Binding(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		return formDecoder.Decode(v, p.Values())
	}

	return ProtocolTypeError
}
