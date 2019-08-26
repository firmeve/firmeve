package queue

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"testing"
	"time"
)

//func TestNewManager(t *testing.T) {
//	z := make(chan int)
//	go func() {
//		z <- 1
//		time.Sleep(time.Second*5)
//		z <- 2
//	}()
//	for {
//		select {
//		case r := <-z:
//			time.Sleep(time.Second*2)
//			fmt.Println(r)
//		case <- time.After(time.Second*1):
//			fmt.Println("超时")
//		}
//
//	}
//	time.Sleep(time.Minute)
//}

type JobDemo struct {
}

func (jd *JobDemo) Handle(payload *Payload) error {
	panic(&Error{Message: "Foo"})
	//panic(`abc`)
	fmt.Println(payload.Data)
	return nil
}
func (jd *JobDemo) Failed(err *Error) {
	//panic(`abc`)
	fmt.Println(err.Payload())
}

func TestQueue(t *testing.T) {

	jobDemo := &JobDemo{}
	manager := NewManager(config.NewConfig("../testdata/config"))

	manager.RegisterJob(`jobDemo`, jobDemo)

	go manager.Run(`default_queue`)

	manager.Connection(`memory`).Push(`jobDemo`, WithQueue(`default_queue`), WithData("abc"))

	time.Sleep(time.Second * 5)
}

var c chan int = make(chan int,0)

//func TestAA(t *testing.T) {
//	go func(c chan int) {
//		fmt.Println(len(c))
//		fmt.Println("--------------")
//		for {
//			if v, ok := <-c; ok {
//				fmt.Println(len(c))
//				fmt.Println("==========")
//				fmt.Println(v)
//			} else {
//				fmt.Println("close")
//			}
//		}
//	}(c)
//	c <- 1
//	close(c)
//	time.Sleep(time.Second * 20)
//	//c <- 1
//	//c <- 2
//	//c <- 3
//	//close(c)
//}
