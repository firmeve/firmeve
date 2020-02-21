package binding

import "github.com/firmeve/firmeve/kernel/contract"

func Protocol(protocol contract.Protocol, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		if p.IsJson() {
			return JSON.Protocol(protocol, v)
		} else {
			//todo
		}
	}

	return nil
}
