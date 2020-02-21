package binding

import "github.com/firmeve/firmeve/kernel/contract"

type (
	multipartForm struct {
	}
)

var (
	MultipartForm = multipartForm{}
)

func (multipartForm) Protocol(protocol contract.Protocol, v interface{}) error {
	//@todo 暂时不解析file
	return Form.Protocol(protocol, v)
}

//@todo 这里后面会操作返回一个file对象类型
func (multipartForm) File() {

}
