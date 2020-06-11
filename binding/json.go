package binding

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

func (json) Protocol(protocol contract.Protocol, v interface{}) error {
	message, _ := protocol.Message()
	return json2.Unmarshal(message, v)
}

//func (json) Data(data []byte, v interface{}) error {
//	return json2.Unmarshal(data, v)
//}
