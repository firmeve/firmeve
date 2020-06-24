package render

import (
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	plain struct {
	}
)

var (
	Plain = plain{}
)

func (plain) Render(protocol contract.Protocol, status int, v interface{}) error {
	return Text.Render(protocol, status, v)
}
