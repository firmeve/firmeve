package resource

type Item struct {
	*baseResource
}

func NewItem(resource interface{}) *Item {
	return &Item{
		baseResource: newBaseResource(resource),
	}
}

func (i *Item) SetFields(fields ...string) *Item {
	i.fields = fields
	return i
}

func (i *Item) SetMeta(meta Meta) *Item {
	i.meta = meta
	return i
}

func (i *Item) SetKey(key string) *Item {
	i.key = key
	return i
}

func (i *Item) Key() string {
	return i.key
}

func (i *Item) Resolve() ResolveMap {
	resolveMap := i.baseResource.Resolve()
	if len(i.meta) > 0 {
		resolveMap[`meta`] = i.meta
	}
	return resolveMap
}
