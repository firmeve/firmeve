package queue

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/utils"
	"time"
)

type queueChan map[string]chan *Payload

type memory struct {
	handleList queueChan
	//delayList    queueChan
	//reservedList queueChan
}

func NewMemory(config *config.Config) *memory {
	return &memory{
		handleList: make(queueChan),
		//delayList:    make(queueChan),
		//reservedList: make(queueChan),
	}
}

func (m *memory) Push(job string, options ...utils.OptionFunc) {
	payloadBlock := createPayload(job, options...)

	m.pushToList(m.handleList, payloadBlock)
	//if payloadBlock.delay.Sub(time.Now()) > time.Second {
	//	m.pushToList(m.delayList, payloadBlock)
	//} else {
	//	m.pushToList(m.handleList, payloadBlock)
	//}
}

func (m *memory) PushRaw(payload *Payload) {
	m.pushToList(m.handleList, payload)
}

func (m *memory) pushToList(queue queueChan, payload *Payload) {
	mu.Lock()
	if _, ok := queue[payload.Queue]; !ok {
		queue[payload.Queue] = make(chan *Payload)
	}
	mu.Unlock()
	queue[payload.Queue] <- payload
}

func (m *memory) Pop(queue string) *Payload {
	if v, ok := <-m.handleList[queue]; ok {
		return v
	}
	return nil
}

func (m *memory) Delay(delay time.Duration, job string, options ...utils.OptionFunc) {
	options = append(options, WithDelay(delay))
	m.Push(job, options...)
}
