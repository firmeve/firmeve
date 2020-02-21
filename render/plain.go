package render

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	plain struct {
	}
)

var (
	Plain = plain{}
)

func (plain) Render(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		p.SetHeader(`Content-Type`, `text/plain`)
	}

	if bytes, ok := v.([]byte); ok {
		_, err := protocol.Write(bytes)
		return err
	}

	return fmt.Errorf("value conversion failed %#v", v)
}
