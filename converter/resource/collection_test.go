package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockCollectionStruct() []struct {
	ID      uint `resource:"id"`
	Title   string
	Content string `resource:"_content"`
} {
	return []struct {
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
}

func TestCollection_Resolve(t *testing.T) {
	structs := MockCollectionStruct()
	//newInterface := make([]interface{}, len(structs))
	//for i, s := range structs {
	//	newInterface[i] = s
	//}
	collection := NewCollection(structs, &Option{
		Fields: Fields{"id"},
	})
	collection.SetMeta(Meta{"a": "a"})
	v := collection.CollectionData()
	assert.Equal(t,v,collection.CollectionData())
	assert.Equal(t, `a`, collection.Meta()[`a`])
	assert.Equal(t, uint(10), v[0][`id`].(uint))
	assert.Equal(t, uint(11), v[1][`id`].(uint))
}
