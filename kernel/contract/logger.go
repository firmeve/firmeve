package contract

import "io"

type (
	Loggable interface {
		Writer(channel string) io.Writer
		With(context ...interface{}) Loggable
		Debug(message string, context ...interface{})
		Info(message string, context ...interface{})
		Warn(message string, context ...interface{})
		Error(message string, context ...interface{})
		Fatal(message string, context ...interface{})
		Panic(message string, context ...interface{})
		//Channel(stack string) Loggable
	}
)
