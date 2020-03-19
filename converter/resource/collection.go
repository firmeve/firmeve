package resource

import (
	"github.com/firmeve/firmeve/kernel/contract"
	reflect2 "reflect"

	"github.com/firmeve/firmeve/support/reflect"
)

type Collection struct {
	resource    []interface{}
	resolveData contract.ResourceDataCollection
	option      *Option
	meta        contract.ResourceMetaData
	link        contract.ResourceLinkData
}

// 还不如,直接baseresource直接解析item,item里面包resource,然后 collection直接再包一层item
// baseresource就不直接解析resource了
// item里面增加一个接口就是resource

func NewCollection(resource interface{}, option *Option) *Collection {
	return &Collection{
		resource:    reflect.SliceInterface(reflect2.ValueOf(resource)),
		option:      option,
		resolveData: make(contract.ResourceDataCollection, 0),
	}
}

func (c *Collection) SetMeta(meta contract.ResourceMetaData) {
	c.meta = meta
}

func (c *Collection) Meta() contract.ResourceMetaData {
	return c.meta
}

func (c *Collection) SetLink(link contract.ResourceLinkData) {
	c.link = link
}

func (c *Collection) Link() contract.ResourceLinkData {
	return c.link
}

func (c *Collection) CollectionData() contract.ResourceDataCollection {
	if len(c.resolveData) > 0 {
		return c.resolveData
	}

	for _, source := range c.resource {
		if v, ok := source.(*Item); ok {
			c.resolveData = append(c.resolveData, v.SetOption(c.option).Data())
		} else {
			//@todo 这里后面操作可能会有问题,collection的原始transformer会被覆盖
			if c.option.Transformer != nil {
				c.option.Transformer = reflect2.New(reflect2.TypeOf(c.option.Transformer).Elem()).Interface().(contract.ResourceTransformer)
				//temp := *c.option.Transformer
				//c.option.Transformer = &temp
			}
			c.resolveData = append(c.resolveData, NewItem(source, c.option).Data())
		}
	}

	return c.resolveData
}
