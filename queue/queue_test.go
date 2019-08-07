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

func (jd *JobDemo) Handle(data interface{}) error {
	panic(&Error{Message:"Foo"})
	//panic(`abc`)
	fmt.Println(data)
	return nil
}

func TestQueue(t *testing.T) {

	jobDemo := &JobDemo{}
	manager := NewManager(config.NewConfig("../testdata/config"))

	manager.RegisterJob(`jobDemo`,jobDemo)

	go manager.Run(`default_queue`)

	manager.Connection(`memory`).Push(`jobDemo`,WithQueueName(`default_queue`),WithData("abc"))

	time.Sleep(time.Second * 5)
}