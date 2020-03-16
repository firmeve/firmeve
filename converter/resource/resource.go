package resource

import "github.com/firmeve/firmeve/converter/transform"

type Option struct {
	Transformer transform.Transformer
	Fields      []string
}
