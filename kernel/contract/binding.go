package contract

type (
	Binding interface {
		Name() string

		Protocol(protocol Protocol, v interface{}) error
	}
)
