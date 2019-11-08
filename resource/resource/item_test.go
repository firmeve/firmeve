package resource

import (
	"encoding/json"
	"fmt"
	"testing"

	resource2 "github.com/firmeve/firmeve/resource"

	"github.com/magiconair/properties/assert"
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

	v := NewItem(m).SetFields(`id`, `title`, `inner_mock`, `_content`).Resolve()
	//fmt.Printf("%#v",v)
	//vs, _ := json.Marshal(v)
	//fmt.Println(string(vs))

	assert.Equal(t, "main title", v[`data`].(ResolveMap)[`title`].(string))
	assert.Equal(t, "content", v[`data`].(ResolveMap)[`_content`].(string))
	assert.Equal(t, uint(10), v[`data`].(ResolveMap)[`id`].(uint))
	assert.Equal(t, subM, v[`data`].(ResolveMap)[`inner_mock`])
}

func TestStruct(t *testing.T) {
	subM := &innerMock{
		IId:     11,
		ITitle:  "inner title",
		content: "inner content",
	}
	m := mock{
		ID:      10,
		Title:   "main title",
		Content: "content",
		SubMock: subM,
	}

	v := NewItem(m).SetFields(`id`, `title`, `inner_mock`, `_content`).Resolve()
	vs, _ := json.Marshal(v)
	fmt.Println(string(vs))

	assert.Equal(t, "main title", v[`data`].(ResolveMap)[`title`].(string))
	assert.Equal(t, "content", v[`data`].(ResolveMap)[`_content`].(string))
	assert.Equal(t, uint(10), v[`data`].(ResolveMap)[`id`].(uint))
	assert.Equal(t, subM, v[`data`].(ResolveMap)[`inner_mock`])
}

type AppTransformer struct {
	resource2.Transformer
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

	v := NewItem(resource2.New(source, &AppTransformer{})).SetFields(`id`, `title`).Resolve()
	fmt.Println(v[`data`].(ResolveMap)[`id`].(uint))
	assert.Equal(t, uint(11), v[`data`].(ResolveMap)[`id`].(uint))
}
