package binding

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"mime/multipart"
)

type (
	multipartForm struct {
		//files map[string]multipart.File
	}
)

var (
	MultipartForm = multipartForm{
		//files: make(map[string]multipart.File, 0),
	}
)

func (multipartForm) Protocol(protocol contract.Protocol, v interface{}) error {
	//@todo 暂时不解析file
	return Form.Protocol(protocol, v)
}

func (m *multipartForm) parseFile() {
}

//@todo 这里后面会操作返回一个file对象类型
func (m *multipartForm) File(protocol contract.Protocol, key string) (multipart.File, *multipart.FileHeader, error) {
	return nil, nil, nil
	//return protocol
	//func (j *MultipartForm) File( {
	//	if j.original.File != nil {
	//		if fhs := j.original.File[key]; len(fhs) > 0 {
	//			f, err := fhs[0].Open()
	//			return f, fhs[0], err
	//		}
	//	}
	//	return nil, nil, http.ErrMissingFile
	//}
}
