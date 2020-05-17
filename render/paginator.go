package render

import (
	"fmt"
	"github.com/firmeve/firmeve/converter/resource"
	"github.com/firmeve/firmeve/converter/serializer"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/ulule/paging"
)

var Paginator = paginator{}

type (
	paginator struct {
	}

	PaginatorResource struct {
		Data        *paging.GORMStore
		Transformer contract.ResourceTransformer
		Fields      contract.ResourceFields
		Meta        contract.ResourceMetaData
		Link        contract.ResourceLinkData
		Limit       uint
	}
)

func (p paginator) Render(protocol contract.Protocol, status int, v interface{}) error {
	if value, ok := v.(*PaginatorResource); ok {
		if value.Limit == 0 {
			value.Limit = 15
		}
		paginator := resource.NewPaginator(value.Data, &resource.Option{
			Transformer: value.Transformer,
			Fields:      value.Fields,
		}, protocol.(contract.HttpProtocol).Request(), &paging.Options{
			DefaultLimit:  int64(value.Limit),
			MaxLimit:      int64(value.Limit + 10),
			LimitKeyName:  "limit",
			OffsetKeyName: "offset",
		})
		if value.Meta != nil {
			paginator.SetMeta(value.Meta)
		}
		if value.Link != nil {
			paginator.SetLink(value.Link)
		}
		return JSON.Render(protocol, status, serializer.NewData(paginator).Resolve())
	}

	return fmt.Errorf("paginator type error %T", v)
}
