package context

import (
	"fmt"
	"testing"
)

type c struct {
	Message []byte

}

func TestNew(t *testing.T) {
	v := new(c)
	//v.Message = make([]byte,0)
	fmt.Printf("%#v",v)
}
