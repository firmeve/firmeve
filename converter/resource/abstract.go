package resource

import (
	"github.com/firmeve/firmeve/converter/transform"
)

type (
	IMeta interface {
		SetMeta(meta Meta)
		Meta() Meta
	}

	ILink interface {
		SetLink(links Link)
		Link() Link
	}

	Datable interface {
		Data() Data
	}

	CollectionData interface {
		CollectionData() DataCollection
	}

	Option struct {
		Transformer transform.Transformer
		Fields      []string
	}

	Meta map[string]interface{}

	Data map[string]interface{}

	DataCollection []Data

	Link map[string]string

	Fields []string
)
