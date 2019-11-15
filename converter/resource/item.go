package resource

type Item struct {
	*baseResource
}

func NewItem(resource interface{}) *Item {
	item := &Item{
		baseResource: newBaseResource(resource),
	}
	//option2 := support.ApplyOption(&option{}, options...)
	//item.transformer = option2.transformer

	return item
}

func (i *Item) SetFields(fields ...string) *Item {
	i.fields = fields
	return i
}

func (i *Item) SetMeta(meta Meta) {
	i.meta = meta
}
func (i *Item) Meta() Meta {
	return i.meta
}

func (i *Item) SetLink(link Link) {
	i.link = link
}
func (i *Item) Link() Link {
	return i.link
}

//func (i *Item) SetKey(key string) *Item {
//	i.key = key
//	return i
//}

//func (i *Item) Key() string {
//	return i.key
//}

func (i *Item) Data() Data {
	return i.baseResource.resolve()
}
