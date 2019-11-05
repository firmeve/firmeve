package resource

type Transformer interface {
	Source() interface{}
}

type BaseTransformer struct {
}

func (t *BaseTransformer) t1Field() string {
	return `abc`
}

func (t *BaseTransformer) t2Field() int {
	return 2
}

func (t *BaseTransformer) t3Field() map[string]string {
	return map[string]string{"a": "a", "b": "b"}
}
