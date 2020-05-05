package contract

import "io"

type (
	// 整个协议的生命周期
	// 包括，协议名称，本身详情，数据完整的数据报，普通value，文件生成，以及协议返回的输出
	Protocol interface {
		io.Reader
		io.Writer

		Application() Application

		Name() string

		// Protocol metadata
		Metadata() map[string][]string

		// Full protocol message
		Message() ([]byte, error)

		Values() map[string][]string
	}
)
