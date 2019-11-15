package serializer

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"fmt"

	"github.com/firmeve/firmeve/converter/resource"
)

func mockCollectionStruct() []struct {
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

func TestNewData_Resolve(t *testing.T) {
	collection := resource.NewCollection(mockCollectionStruct()).SetFields(`id`, `title`)
	collection.SetMeta(resource.Meta{"a": 1, "head": "head"})
	//fmt.Printf("%#v", collection.CollectionData())
	v := NewData(collection).Resolve()

	b, _ := xml.Marshal(collection)
	fmt.Println(string(b))

	fmt.Printf("%#v\n", v)
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bytes)
}
