package utils

type OptionFunc func(option interface{})

// Apply option
func ApplyOption(option interface{}, options ...OptionFunc) {
	for _, opt := range options {
		opt(option)
	}
}
