package contract

type (
	Provider interface {
		Name() string
		Register()
		Boot()
	}
)
