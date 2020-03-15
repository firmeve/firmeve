package contract

type (
	Render interface {
		Render(protocol Protocol, status int, v interface{}) error
	}
)
