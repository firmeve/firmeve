package contract

type (
	Loggable interface {
		Debug(message string, context ...interface{})
		Info(message string, context ...interface{})
		Warn(message string, context ...interface{})
		Error(message string, context ...interface{})
		Fatal(message string, context ...interface{})
		Panic(message string, context ...interface{})
		//Channel(stack string) Loggable
	}
)
