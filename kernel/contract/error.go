package contract

type Error interface {
	error

	Details() []interface{}

	Equal(err Error) bool

	Unwrap() error
}
