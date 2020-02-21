package contract

type (
	Render interface {
		Render(protocol Protocol, v interface{}) error
	}
)
