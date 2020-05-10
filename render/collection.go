package render

import (
	"fmt"
	"github.com/firmeve/firmeve/converter/resource"
	"github.com/firmeve/firmeve/kernel/contract"
)

var Collection = collection{}

type (
	collection struct {
	}

	CollectionResource struct {
		Data   interface{}
		Option *resource.Option
		//Meta   contract.ResourceMetaData
		//Link   contract.ResourceLinkData
	}
)

func (c collection) Render(protocol contract.Protocol, status int, v interface{}) error {
	if value, ok := v.(*CollectionResource); ok {
		return JSON.Render(protocol, status, resource.NewCollection(value.Data, value.Option))
	}

	return fmt.Errorf("collection type error %T", v)
}
