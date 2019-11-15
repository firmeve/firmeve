package resource

import (
	reflect2 "reflect"

	"github.com/firmeve/firmeve/support/reflect"
)

type Collection struct {
	*baseResource
	resource []interface{}
}

func NewCollection(resource interface{}) *Collection {
	return &Collection{
		resource:     reflect.SliceInterface(reflect2.ValueOf(resource)),
		baseResource: newBaseResource(nil),
	}
}

func (c *Collection) SetFields(fields ...string) *Collection {
	c.fields = fields
	return c
}

func (c *Collection) SetMeta(meta Meta) {
	c.meta = meta
}

func (c *Collection) Meta() Meta {
	return c.meta
}

func (c *Collection) SetLink(link Link) {
	c.link = link
}
func (c *Collection) Link() Link {
	return c.link
}

func (c *Collection) CollectionData() DataCollection {
	dataMaps := make(DataCollection, 0)
	for _, source := range c.resource {
		if v, ok := source.(*Item); ok {
			dataMaps = append(dataMaps, v.SetFields(c.fields...).Data())
		} else {
			dataMaps = append(dataMaps, NewItem(source).SetFields(c.fields...).Data())
		}
	}

	return dataMaps
}
