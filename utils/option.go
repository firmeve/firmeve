package utils

type Option interface{}

type OptionFunc func(option Option)

// Apply option
func ApplyOption(option Option, options ...OptionFunc) Option {
	for _, opt := range options {
		opt(option)
	}

	return option
}
