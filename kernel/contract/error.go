package contract

type Error interface {
	error

	Details() []interface{}

	Equal(err error) bool

	Unwrap() error

	String() string
}
