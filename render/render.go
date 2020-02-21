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

func Render(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		//@todo accept-type可能有多个
		acceptType := p.Header(`Accept-Type`)
		if r, ok := httpRenderType[acceptType]; ok {
			return r.Render(protocol, v)
		}

		return fmt.Errorf("non-existent type %s", acceptType)
	}

	return nil
}
