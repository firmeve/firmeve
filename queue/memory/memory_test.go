package memory

import (
	"fmt"
	"testing"
)

type JobFunc func(interface{})

func (f JobFunc) Handle(data interface{}) {
	f(data)
}

func TestMemory(t *testing.T) {
	queue := NewMemory()

	go queue.Run()

	queue.Push(JobFunc(func(i interface{}) {
		fmt.Println(i)
	}),"sss","fdsfdsfds")

	//http.HandleFunc()

	//JobFunc(func(i interface{}) {
	//	fmt.Println(i)
	//}).Handle("abc")


	//time.Sleep(time.Hour)
}
