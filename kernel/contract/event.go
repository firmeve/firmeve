package contract

type (
	EventHandler interface {
		Handle(application Application, params ...interface{}) (interface{}, error)
	}

	Event interface {
		Listen(name string, handlers ...EventHandler)
		Dispatch(name string, params ...interface{}) []interface{}
		Has(name string) bool
	}
)
