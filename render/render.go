package render

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
)

var (
	httpRenderType = map[string]contract.Render{
		contract.HttpMimeJson:  JSON,
		contract.HttpMimePlain: Plain,
	}
)

func Render(protocol contract.Protocol, status int, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		accept := p.Accept()

		for _, item := range accept {
			if r, ok := httpRenderType[item]; ok {
				return r.Render(protocol, status, v)
			}
		}

		return fmt.Errorf("non-existent type %v", accept)
	}

	return nil
}
