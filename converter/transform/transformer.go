package transform

type Transformer interface {
	Resource() interface{}
	SetResource(resource interface{})
}

type BaseTransformer struct {
	Original interface{}
}

func (t *BaseTransformer) Resource() interface{} {
	return t.Original
}

func (t *BaseTransformer) SetResource(resource interface{}) {
	t.Original = resource
}
