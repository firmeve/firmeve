package transform

type Transformer interface {
	SetResource(resource interface{})
	Resource() interface{}
}

type BaseTransformer struct {
	resource interface{}
}

func (t *BaseTransformer) Resource() interface{} {
	return t.resource
}

func (t *BaseTransformer) SetResource(resource interface{}) {
	t.resource = resource
}
