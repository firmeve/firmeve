package scheduler

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"sync"
	"testing"
	"time"
)

var s2 = New(&Configuration{
	Available: 2,
})
var once2 = sync.Once{}

type handler struct {
}

func (h handler) Handle(message *contract.SchedulerMessage) error {
	time.Sleep(time.Second)
	fmt.Println(message.Worker, message.Message)
	return nil
}

func (h handler) Failed(err error) {
	panic("implement me")
}

func TestNewScheduler(t *testing.T) {
	once2.Do(func() {
		s2.RegisterHandler(`testing`, new(handler))
		s2.Run()
	})

	for i := 0; i < 100; i++ {
		if i == 90 {
			s2.Close()
			break
		} else if i == 10 {
			s2.Increment(3)
		} else if i == 30 {
			s2.Decrement(3)
		} else if i == 45 {
			s2.Decrement(1)
		} else if i == 60 {
			s2.Increment(2)
		}

		s2.Send(&contract.SchedulerMessage{
			Message: fmt.Sprintf(`send message %d`, i),
			Handler: `testing`,
		})
	}
}
