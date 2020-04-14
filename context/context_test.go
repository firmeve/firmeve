package context

import (
	"fmt"
	"testing"
)

type c struct {
	Message []byte
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func TestNew2(t *testing.T) {
	fmt.Println(filterFlags("abc;def"))
	//v := new(c)
	//v.Message = make([]byte,0)
	//fmt.Printf("%#v", v)
}
