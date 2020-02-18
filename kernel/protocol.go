package kernel

import (
	"encoding/json"
	"github.com/firmeve/firmeve/kernel/contract"
)

type Json func(protocol contract.Protocol, v interface{})

func (j *Json) Unpack(protocol contract.Protocol) map[string][]string {
	v := make([]byte,0)
	json.Marshal()
	json.Unmarshal(v,protocol.Message())
	return make(map[string][]string,0)
}

func () Pack(protocol contract.Protocol) {
	panic("implement me")
}

func () File(protocol contract.Protocol) {
	panic("implement me")
}

func () Output(protocol contract.Protocol, v interface{}) []byte {
	panic("implement me")
}
