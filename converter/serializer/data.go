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
	} else {
		collection[`data`] = d.resource
	}

	if meta, ok := d.resource.(resource.IMeta); ok {
		metaData := meta.Meta()
		if len(metaData) > 0 {
			collection[`meta`] = metaData
		}
	}

	if link, ok := d.resource.(resource.ILink); ok {
		linkData := link.Link()
		if len(linkData) > 0 {
			collection[`link`] = linkData
		}
	}

	return collection
}

func NewData(resource interface{}) *Data {
	return &Data{resource: resource}
}
