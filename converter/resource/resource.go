package resource

import (
	"github.com/firmeve/firmeve/kernel/contract"
)

type Option struct {
	Transformer contract.ResourceTransformer
	Fields      []string
}
