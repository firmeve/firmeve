package scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestNew(t *testing.T) {
	s := New(&Configuration{Size: 10})
	s.RegisterHandler(`handler`, &TestHandler{})
	s.Dispatch()
	//time.Sleep(time.Second)
	//s.Delivery(`h1`, `abc`)
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			rand.Seed(time.Now().UnixNano())
			s.Delivery(`handler`, rand.Int())
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(time.Second * 120)
}

type TestHandler struct {
}

func (h *TestHandler) Handle(ctx context.Context, data interface{}) {
	for {

		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				fmt.Println(ctx.Err())
			}
		default:
			fmt.Println(data)
			time.Sleep(time.Second * 3)
		}

	}

}

func (h *TestHandler) ParentCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	return ctx
}
