package contract

type (
	Validator interface {
		Validate(val interface{}) error
	}
)
