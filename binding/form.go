package binding

import (
	"errors"

	"github.com/firmeve/firmeve/kernel/contract"
	form2 "github.com/go-playground/form/v4"
)

type (
	form struct {
	}
)

var (
	ProtocolTypeError = errors.New("the protocol only support http protocol")
	formDecoder       = form2.NewDecoder()
	Form              = form{}
)

func (f form) Name() string {
	return `x-www-form-urlencoded`
}

func (f form) Protocol(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		return formDecoder.Decode(v, p.Values())
	}

	return ProtocolTypeError
}
