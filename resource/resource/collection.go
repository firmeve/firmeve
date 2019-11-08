package resource

import (
	"github.com/firmeve/firmeve/support/reflect"
	reflect2 "reflect"
)

type Collection struct {
	*baseResource
	resource []interface{}
}

func (c *Collection) SetFields(fields ...string) *Collection {
	c.fields = fields
	return c
}

func (c *Collection) SetMeta(meta Meta) *Collection {
	c.meta = meta
	return c
}

func (c *Collection) SetKey(key string) *Collection {
	c.key = key
	return c
}

func (c *Collection) Key() string {
	return c.key
}

func NewCollection(resource interface{}) *Collection {
	return &Collection{
		resource:     reflect.SliceInterface(reflect2.ValueOf(resource)),
		baseResource: newBaseResource(nil),
	}
}

func (c *Collection) Resolve() ResolveMap {
	resolveMaps := make([]ResolveMap, 0)
	for _, source := range c.resource {
		resolveMap := NewItem(source).SetFields(c.fields...).SetKey(``).Resolve()
		resolveMaps = append(resolveMaps, resolveMap)
	}

	collection := make(ResolveMap, 0)
	if c.key == `` {
		panic(`collection key not empty`)
	}
	collection[c.key] = resolveMaps
	if len(c.meta) > 0 {
		collection[`meta`] = c.meta
	}

	return collection
}
