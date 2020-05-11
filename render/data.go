package render

import (
	"github.com/firmeve/firmeve/converter/serializer"
	"github.com/firmeve/firmeve/kernel/contract"
)

var Data = data{}

type (
	data struct {
	}
)

func (data) Render(protocol contract.Protocol, status int, v interface{}) error {
	return JSON.Render(protocol, status, serializer.NewData(v).Resolve())
}
