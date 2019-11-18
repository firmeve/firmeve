package resource

import (
	reflect2 "reflect"

	"github.com/firmeve/firmeve/converter/transform"
	"github.com/firmeve/firmeve/support/reflect"
)

type Collection struct {
	resource []interface{}
	option   *Option
	meta     Meta
	link     Link
}

// 还不如,直接baseresource直接解析item,item里面包resource,然后 collection直接再包一层item
// baseresource就不直接解析resource了
// item里面增加一个接口就是resource

func NewCollection(resource interface{}, option *Option) *Collection {
	return &Collection{
		resource: reflect.SliceInterface(reflect2.ValueOf(resource)),
		option:   option,
	}
}

func (c *Collection) SetMeta(meta Meta) {
	c.meta = meta
}

func (c *Collection) Meta() Meta {
	return c.meta
}

func (c *Collection) SetLink(link Link) {
	c.link = link
}

func (c *Collection) Link() Link {
	return c.link
}

func (c *Collection) CollectionData() DataCollection {
	dataMaps := make(DataCollection, 0)
	for _, source := range c.resource {
		if v, ok := source.(*Item); ok {
			dataMaps = append(dataMaps, v.SetOption(c.option).Data())
		} else {
			//@todo 这里后面操作可能会有问题,collection的原始transformer会被覆盖
			if c.option.Transformer != nil {
				c.option.Transformer = reflect2.New(reflect2.TypeOf(c.option.Transformer).Elem()).Interface().(transform.Transformer)
				//temp := *c.option.Transformer
				//c.option.Transformer = &temp
			}
			dataMaps = append(dataMaps, NewItem(source, c.option).Data())
		}
	}

	return dataMaps
}
