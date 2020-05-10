package render

import (
	"fmt"
	"github.com/firmeve/firmeve/converter/resource"
	"github.com/firmeve/firmeve/converter/serializer"
	"github.com/firmeve/firmeve/kernel/contract"
)

var Item = item{}

type (
	item struct {
	}

	ItemResource struct {
		Data        interface{}
		Transformer contract.ResourceTransformer
		Fields      contract.ResourceFields
		Meta        contract.ResourceMetaData
		Link        contract.ResourceLinkData
	}
)

func (c item) Render(protocol contract.Protocol, status int, v interface{}) error {
	if value, ok := v.(*ItemResource); ok {
		item := resource.NewItem(value.Data, &resource.Option{
			Transformer: value.Transformer,
			Fields:      value.Fields,
		})
		if value.Meta != nil {
			item.SetMeta(value.Meta)
		}
		if value.Link != nil {
			item.SetLink(value.Link)
		}
		return JSON.Render(protocol, status, serializer.NewData(item).Resolve())
	}

	return fmt.Errorf("item type error %T", v)
}
