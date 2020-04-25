package contract

type (
	EventHandler interface {
		Handle(params ...interface{}) (interface{}, error)
	}

	Event interface {
		Listen(name string, handler EventHandler)
		ListenMany(name string, handlers []EventHandler)
		Dispatch(name string, params ...interface{}) []interface{}
		Has(name string) bool
	}
)
