package contract

type (
	Binding interface {
		Protocol(protocol Protocol, v interface{}) error
	}
)
