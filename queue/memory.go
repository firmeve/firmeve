package queue

import (
	"fmt"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/utils"
	"time"
)

type queueChan map[string]chan *Payload

type memory struct {
	handleList   queueChan
	delayList    queueChan
	reservedList queueChan
}

func NewMemory(config *config.Config) *memory {
	return &memory{
		handleList:   make(queueChan),
		delayList:    make(queueChan),
		reservedList: make(queueChan),
	}
}

func (m *memory) Push(jobName string, options ...utils.OptionFunc) {
	payloadBlock := createPayload(jobName, options...)

	if payloadBlock.delay.Sub(time.Now()) > time.Second {
		m.pushToList(m.delayList, payloadBlock)
	} else {
		m.pushToList(m.handleList, payloadBlock)
	}
}

func (m *memory) pushToList(queue queueChan, payload *Payload) {
	mu.Lock()
	if _, ok := queue[payload.queueName]; !ok {
		queue[payload.queueName] = make(chan *Payload)
	}
	mu.Unlock()
	queue[payload.queueName] <- payload
}

func (m *memory) Pop(queueName string) *Payload {
	//mu.Lock()
	//defer mu.Unlock()
	if v, closed := <- m.handleList[queueName]; closed {
		fmt.Println(v)
		return v
	}
	return nil
	//return m.handleList[queueName]
}

func (m *memory) Delay(delay time.Time, jobName string, options ...utils.OptionFunc) {
	options = append(options, WithDelay(delay))
	m.Push(jobName, options...)
}
