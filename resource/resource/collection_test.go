package resource

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollection_Resolve(t *testing.T) {
	structs := []struct {
		ID      uint `resource:"id"`
		Title   string
		Content string `resource:"_content"`
	}{
		{
			ID:      10,
			Title:   "main title",
			Content: "content",
		},
		{
			ID:      11,
			Title:   "main title1",
			Content: "content1",
		},
		{
			ID:      12,
			Title:   "main title2",
			Content: "content2",
		},
	}
	newInterface := make([]interface{}, len(structs))
	for i, s := range structs {
		newInterface[i] = s
	}
	v := NewCollection(structs).SetFields(`id`, `title`).SetMeta(Meta{"a": "a"}).Resolve()
	assert.Equal(t, `a`, v[`meta`].(Meta)[`a`])
	assert.Equal(t, uint(10), v[`data`].([]ResolveMap)[0][`id`].(uint))
	assert.Equal(t, uint(11), v[`data`].([]ResolveMap)[1][`id`].(uint))
}
