package demo

type T1 struct {
	Name string
}

func NewT1() *T1 {
	return &T1{"Simon"}
}

type T2 struct {
	t1 *T1
	Age int
}

func NewT2(f *T1) T2  {
	return T2{t1:f,Age:10}
}
