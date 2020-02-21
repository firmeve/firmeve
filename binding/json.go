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

func (json) Name() string {
	return `json`
}

func (json) Protocol(protocol contract.Protocol, v interface{}) error {
	//if p, ok := protocol.(contract.HttpProtocol); ok && p.IsJson() {
	//}
	//if protocol.Name() == `http` && protocol.Metadata() {
	//
	//}
	message, _ := protocol.Message()
	return json2.Unmarshal(message, v)
}

func (json) Data(data []byte, v interface{}) error {
	return json2.Unmarshal(data, v)
}

//type (
//	JSONData map[string]interface{}
//
//	JSON struct {
//		original []byte
//		data     JSONData
//	}
//)
//
//func (j *JSON) Name() string {
//	panic("implement me")
//}
//
//func (j *JSON) Binding(protocol contract.Protocol, v interface{}) error {
//	panic("implement me")
//}
//
//func (j *JSON) Bind(object interface{}) error {
//	return json.Unmarshal(j.original, object)
//}
//
//func (j *JSON) Has(key string) bool {
//	_, ok := j.data[key]
//	return ok
//}
//
//func (j *JSON) Get(key string) interface{} {
//	if j.Has(key) {
//		return j.data[key]
//	}
//
//	return ``
//}
//
//func (j *JSON) parse() {
//	if err := json.Unmarshal(j.original, &j.data); err != nil {
//		panic(err)
//	}
//}
//
//func NewJSON(data []byte) *JSON {
//	j := &JSON{
//		original: data,
//		data:     make(JSONData, 0),
//	}
//	j.parse()
//	return j
//}
