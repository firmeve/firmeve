package contract

type (
	Binding interface {
		Name() string

		Binding(protocol Protocol, v interface{}) error
	}
)
