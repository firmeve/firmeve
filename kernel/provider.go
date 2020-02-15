package kernel

import "github.com/firmeve/firmeve/kernel/contract"

type (
	BaseProvider struct {
		Firmeve contract.Application `inject:"firmeve"`
	}
)
