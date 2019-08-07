package queue

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"testing"
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

func (jd *JobDemo) Handle(data interface{}) {
	fmt.Println(data)
}

func TestQueue(t *testing.T) {
	jobDemo := &JobDemo{}

	RegisterJob(`jobDemo`,jobDemo)
	go Run(`default_queue`)
	manager := NewManager(config.NewConfig("../testdata/config"))
	manager.Connection(`memory`).Push(`jobDemo`,WithQueueName(`default_queue`),WithData("abc"))


}