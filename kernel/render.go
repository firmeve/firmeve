package kernel

import "github.com/firmeve/firmeve/kernel/contract"

func Render(protocol contract.Protocol,v interface{})  {
	switch protocol.Name() {
	case `http`:
		//content-type
		protocol.Write()
	}
}
