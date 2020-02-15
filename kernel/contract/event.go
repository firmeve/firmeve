package contract

type (
	EventHandler interface {
		Handle(params map[string]interface{}) (interface{}, error)
	}

	Event interface {
		Listen(name string, handler EventHandler)
		ListenMany(name string, handlers []EventHandler)
		Dispatch(name string, params map[string]interface{}) []interface{}
		Has(name string) bool
	}
)