package resource

import (
	"testing"

	"github.com/firmeve/firmeve/converter/transform"

	"github.com/stretchr/testify/assert"
)

type innerMock struct {
	IId     uint
	ITitle  string
	content string
}

type mock struct {
	ID      uint `resource:"id"`
	Title   string
	Content string     `resource:"_content"`
	SubMock *innerMock `resource:"inner_mock"`
}

func TestPrtStruct(t *testing.T) {
	subM := &innerMock{
		IId:     11,
		ITitle:  "inner title",
		content: "inner content",
	}
	m := &mock{
		ID:      10,
		Title:   "main title",
		Content: "content",
		SubMock: subM,
	}

	v := NewItem(m, &Option{
		Fields: []string{`id`, `title`, `inner_mock`, `_content`},
	}).Data()
	//fmt.Println(v)
	//fmt.Printf("%#v", v)
	//vs, _ := json.Marshal(v)
	//fmt.Println(string(vs))

	assert.Equal(t, "main title", v[`title`].(string))
	assert.Equal(t, "content", v[`_content`].(string))
	assert.Equal(t, uint(10), v[`id`].(uint))
	assert.Equal(t, subM, v[`inner_mock`])
}

//func TestStruct(t *testing.T) {
//	subM := &innerMock{
//		IId:     11,
//		ITitle:  "inner title",
//		content: "inner content",
//	}
//	m := mock{
//		ID:      10,
//		Title:   "main title",
//		Content: "content",
//		SubMock: subM,
//	}
//
//	v := NewItem(m).SetFields(`id`, `title`, `inner_mock`, `_content`).resolve()
//	vs, _ := json.Marshal(v)
//	fmt.Println(string(vs))
//
//	assert.Equal(t, "main title", v[`title`].(string))
//	assert.Equal(t, "content", v[`_content`].(string))
//	assert.Equal(t, uint(10), v[`id`].(uint))
//	assert.Equal(t, subM, v[`inner_mock`])
//}

type AppTransformer struct {
	transform.BaseTransformer
}

func (t *AppTransformer) IdField() uint {
	return t.Resource().(Source).ID + 1
}

type Source struct {
	ID      uint `resource:"id"`
	Title   string
	Content string `resource:"_content"`
}

func TestTransformer(t *testing.T) {
	source := Source{
		ID:      10,
		Title:   "main title",
		Content: "content",
	}

	item := NewItem(source, &Option{
		Transformer: new(AppTransformer),
		Fields: []string{
			`id`, `title`,
		},
	})
	item.SetMeta(map[string]interface{}{"a": 1})
	v := item.Data()
	m := item.Meta()
	assert.Equal(t, uint(11), v[`id`].(uint))
	assert.Equal(t, 1, m[`a`])
}
