package kernel

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	//v := Error("abcdef")
	//f()

	v := Error("abc")
	fmt.Println(v.StackString())

	//fmt.Printf(Errorf("a %w a", errors.New("nnnn")))

	//e := b()
	//s := WithError("dddd", e.(contract.Error))
	//var b strings.Builder
	//b.WriteString("==========")
	//b.WriteString("\n\n")
	//b.WriteString("Traceback:")
	//for _, pc := range s.Stack() {
	//	fn := runtime.FuncForPC(pc)
	//	b.WriteString("\n")
	//	f, n := fn.FileLine(pc)
	//	b.WriteString(fmt.Sprintf("%s:%d", f, n))
	//}
	//fmt.Println(b.String())

	//e := &Error{
	//	message: "abc",
	//}
	//
	//fmt.Printf("a%wa", errors.New("nnnn"))
}

func f() {
	v := Error("abc")
	fmt.Println(v.StackString())
}
