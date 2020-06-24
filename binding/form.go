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
	Form              = form{}
	decoder           = form2.NewDecoder()
)

func (f form) Protocol(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		return decoder.Decode(v, p.Values())
	}

	return ProtocolTypeError
}
