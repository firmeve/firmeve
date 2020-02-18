package contract

//import "errors"
type Error2 interface {
	error
	//Error() string

	Code() int

	Details() []interface{}

	Equal(err error) bool
}

//func c()  {
//	errors.Unwrap()
//}