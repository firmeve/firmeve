package resource

type Collection struct {
	*Resource
}

//func (c *Collection) Resolve() {
//	collection := make([]ResolveMap, 0)
//	c.Resource.Resolve()
//}
//
//func NewCollection(sources []interface{}) []*Item {
//	collection := make([]*Item, 0)
//	for _, source := range sources {
//		collection = append(collection, NewItem(source))
//	}
//	return collection
//}
