package render

import (
	"fmt"
	"github.com/firmeve/firmeve/converter/resource"
	"github.com/firmeve/firmeve/converter/serializer"
	"github.com/firmeve/firmeve/kernel/contract"
)

var Collection = collection{}

type (
	collection struct {
	}

	CollectionResource struct {
		Data        interface{}
		Transformer contract.ResourceTransformer
		Fields      contract.ResourceFields
		Meta        contract.ResourceMetaData
		Link        contract.ResourceLinkData
	}
)

func (c collection) Render(protocol contract.Protocol, status int, v interface{}) error {
	if value, ok := v.(*CollectionResource); ok {
		collect := resource.NewCollection(value.Data, &resource.Option{
			Transformer: value.Transformer,
			Fields:      value.Fields,
		})
		if value.Meta != nil {
			collect.SetMeta(value.Meta)
		}
		if value.Link != nil {
			collect.SetLink(value.Link)
		}
		return JSON.Render(protocol, status, serializer.NewData(collect).Resolve())
	}

	return fmt.Errorf("collection type error %T", v)
}
