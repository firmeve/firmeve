package support

// Object of operation
type Object interface{}

// Returned object set option function
type Option func(object Object)

// Apply option
func ApplyOption(object Object, options ...Option) Object {
	for _, option := range options {
		option(object)
	}

	return object
}
