package resource

type Resource interface {
	Resource() interface{}
	SetResource(resource interface{})
}

type Transformer struct {
	resource interface{}
}

//
func (t *Transformer) Resource() interface{} {
	return t.resource
}

func (t *Transformer) SetResource(resource interface{}) {
	t.resource = resource
}

func New(resource interface{}, transformer Resource) Resource {
	transformer.SetResource(resource)
	return transformer
}
