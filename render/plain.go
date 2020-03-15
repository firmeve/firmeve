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

func (plain) Render(protocol contract.Protocol, status int, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		p.ResponseWriter().WriteHeader(status)
		p.SetHeader(`Content-Type`, `text/plain`)
	}

	if bytes, ok := v.([]byte); ok {
		_, err := protocol.Write(bytes)
		return err
	}

	return fmt.Errorf("value conversion failed %#v", v)
}
