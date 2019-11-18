package resource

import (
	"github.com/firmeve/firmeve/converter/transform"
)

type Meta map[string]interface{}
type Data map[string]interface{}
type DataCollection []Data
type Link map[string]string
type Fields []string

type IMeta interface {
	SetMeta(meta Meta)
	Meta() Meta
}

type ILink interface {
	SetLink(links Link)
	Link() Link
}

type Datable interface {
	Data() Data
}
type CollectionData interface {
	CollectionData() DataCollection
}

type Option struct {
	Transformer transform.Transformer
	Fields      []string
}
