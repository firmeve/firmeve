package contract

type (
	Render interface {
		Name() string
		Render(protocol Protocol,v interface{}) (error)
	}
)
