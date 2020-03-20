package render

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	html struct {
	}
)

var (
	Html = html{}
)

func (html) Render(protocol contract.Protocol, status int, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		p.ResponseWriter().WriteHeader(status)
		p.SetHeader(`Content-Type`, `text/html`)
	}

	var err error

	_, err = protocol.Write([]byte(fmt.Sprintf("%v", v)))
	//if bytes, ok := v.([]byte); ok {
	//	_, err = protocol.Write(bytes)
	//} else {
	//	_, err = protocol.Write([]byte(fmt.Sprintf("%v", v)))
	//}

	return err
	//return fmt.Errorf("value conversion failed %#v", v)
}
