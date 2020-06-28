package socket

import (
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	Socket struct {
	}
)

func (s Socket) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (s Socket) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func (s Socket) Application() contract.Application {
	panic("implement me")
}

func (s Socket) Name() string {
	panic("implement me")
}

func (s Socket) Metadata() map[string][]string {
	panic("implement me")
}

func (s Socket) Message() ([]byte, error) {
	panic("implement me")
}

func (s Socket) Values() map[string][]string {
	panic("implement me")
}

func (s Socket) Clone() contract.Protocol {
	panic("implement me")
}
