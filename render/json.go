package render

import (
	json2 "encoding/json"
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	json struct {
	}
)

var (
	JSON = json{}
)

func (json) Render(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		p.SetHeader(`Content-Type`, `application/json`)
	}

	bytes, err := json2.Marshal(v)
	if err != nil {
		return err
	}
	_, err = protocol.Write(bytes)
	return err
}
