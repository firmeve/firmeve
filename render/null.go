package render

import "github.com/firmeve/firmeve/kernel/contract"

type null struct {
}

var Null = null{}

func (n null) Render(protocol contract.Protocol, status int, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		p.ResponseWriter().WriteHeader(status)
	}

	_, err := protocol.Write(nil)

	return err
}
