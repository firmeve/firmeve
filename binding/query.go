package binding

import "github.com/firmeve/firmeve/kernel/contract"

type (
	query struct {
	}
)

var (
	Query = query{}
)

func (query) Name() string {
	return `url-query`
}

func (query) Protocol(protocol contract.Protocol, v interface{}) error {
	panic("not support")
}
