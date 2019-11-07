package testdata

type T1 struct {
	Name string
}

func NewT1() *T1 {
	return &T1{"Simon"}
}

func NewT1Error() (*T1, error) {
	return &T1{"Simon"}, nil
}

func NewT1Sturct() T1 {
	return T1{"Simon"}
}

type T2 struct {
	t1  *T1
	S1  *T1 `inject:"t1"`
	TS  T1
	Age int
}

func (t *T2) GetT1() *T1 {
	return t.t1
}

func NewT2(f *T1) T2 {
	return T2{t1: f, Age: 10}
}
func NewT2Error(f *T1, err error) (T2, error) {
	return T2{t1: f, Age: 10}, err
}

func NewT2ErrorInterface(f *T1, err error) (T2, error) {
	//var err error
	//if len(params) == 1 {
	//	err = nil
	//} else {
	//	err = params[1].(error)
	//}
	return T2{t1: f, Age: 10}, err
}

func NewT22() *T2 {
	return &T2{Age: 10}
}

func NewTStruct(t1 T1) *T2 {
	return &T2{Age: 10, TS: t1}
}

func T2Call(f *T1) *T1 {
	//return T2{t1:f,Age:10}
	return f
}
