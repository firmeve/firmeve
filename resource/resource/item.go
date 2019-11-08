package resource

import "fmt"

type Item struct {
	*Resource
}

func NewItem(source interface{}) *Item {
	return &Item{
		Resource: New(source),
	}
}

func (i *Item) Resolve() ResolveMap {
	fmt.Println("111")
	return i.Resource.Resolve()
}
