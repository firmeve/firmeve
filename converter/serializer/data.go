package serializer

import (
	"github.com/firmeve/firmeve/kernel/contract"
)

type Data struct {
	resource interface{}
}

func (d *Data) Resolve() interface{} {
	collection := make(ResolveData, 0)

	if data, ok := d.resource.(contract.ResourceDatable); ok {
		collection[`data`] = data.Data()
	} else if data, ok := d.resource.(contract.ResourceCollectionData); ok {
		collection[`data`] = data.CollectionData()
	} else {
		collection[`data`] = d.resource
	}

	if meta, ok := d.resource.(contract.ResourceMeta); ok {
		metaData := meta.Meta()
		if len(metaData) > 0 {
			collection[`meta`] = metaData
		}
	}

	if link, ok := d.resource.(contract.ResourceLink); ok {
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
