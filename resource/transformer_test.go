package resource

type TestTransformer struct {
	Transformer
	resource interface{}
}

type Source struct {
	ID      uint `resource:"id"`
	Title   string
	Content string `resource:"_content"`
}

//func TestNew(t *testing.T) {
//
//}
