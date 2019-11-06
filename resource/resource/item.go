package resource

type Item struct {
	*Resource
}

func NewItem(source interface{}) *Item {
	return &Item{
		Resource: New(source),
	}
}
