package render

import (
	"fmt"
	"github.com/firmeve/firmeve/converter/resource"
	"github.com/firmeve/firmeve/kernel/contract"
)

var Item = item{}

type (
	item struct {
	}

	ItemResource struct {
		Data   interface{}
		Option *resource.Option
		//Meta   contract.ResourceMetaData
		//Link   contract.ResourceLinkData
	}
)

func (c item) Render(protocol contract.Protocol, status int, v interface{}) error {
	if value, ok := v.(*ItemResource); ok {
		return JSON.Render(protocol, status, resource.NewItem(value.Data, value.Option))
	}

	return fmt.Errorf("item type error %T", v)
}
