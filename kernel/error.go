package kernel

type (
	Error struct {
		err error
		code int
		message string
	}
)

func (e *Error) Error() string {
	panic("implement me")
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Details() []interface{} {
	panic("implement me")
}

func (e *Error) Equal(err error) bool {
}


func Error2()  {
	
}

func Errorf()  {
	
}