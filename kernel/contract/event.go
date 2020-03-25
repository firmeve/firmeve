package contract

type (
	EventInParams map[string]interface{}

	EventHandler interface {
		Handle(params EventInParams) (interface{}, error)
	}

	Event interface {
		Listen(name string, handler EventHandler)
		ListenMany(name string, handlers []EventHandler)
		Dispatch(name string, params EventInParams) []interface{}
		Has(name string) bool
	}
)
