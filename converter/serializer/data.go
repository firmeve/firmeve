package serializer

import (
	"github.com/firmeve/firmeve/converter/resource"
)

type Data struct {
	resource interface{}
}

func (d *Data) Resolve() interface{} {
	collection := make(ResolveData, 0)

	if data, ok := d.resource.(resource.Datable); ok {
		collection[`data`] = data.Data()
	} else if data, ok := d.resource.(resource.CollectionData); ok {
		collection[`data`] = data.CollectionData()
	}

	if meta, ok := d.resource.(resource.IMeta); ok {
		collection[`meta`] = meta.Meta()
	}
	if link, ok := d.resource.(resource.ILink); ok {
		collection[`link`] = link.Link()
	}

	return collection
}

func NewData(resource interface{}) *Data {
	return &Data{resource: resource}
}
