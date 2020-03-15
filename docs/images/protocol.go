package context

import "io"

type (
	//Request interface {
	//	// @todo file struct or file interface
	//	// @todo Unified package
	//	//File()
	//
	//	// @todo stream
	//	//Stream()
	//
	//	// Protocol content
	//	Message() []byte
	//
	//	//
	//	Value(key string) interface{}
	//}
	//
	//Response interface {
	//	//Code(code int)
	//	Write([]byte) (int,error)
	//}

	// 整个协议的生命周期
	// 包括，协议名称，本身详情，数据完整的数据报，普通value，文件生成，以及协议返回的输出
	Protocol interface {
		//Request
		//Response
		// Protocol name
		//io.ReadWriteCloser
		//io.Reader
		//io.ReaderAt
		//io.Seeker,
		io.Reader
		io.Writer
		//io.Seeker

		Name() string

		// Protocol metadata
		Metadata() map[string][]string

		// @todo file struct or file interface
		// @todo Unified package
		//File()

		// @todo stream
		//Stream()

		// Full protocol message
		Message() ([]byte,error)

		Status(code int)

		//
		Values(key string) map[string][]string



		//Write([]byte) (int,error)
	}
)