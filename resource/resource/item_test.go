package resource

import (
	"fmt"
	"reflect"
	"testing"
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

func TestNewItem(t *testing.T) {
	m := &mock{
		ID:      10,
		Title:   "main title",
		Content: "content",
		SubMock: &innerMock{
			IId:     11,
			ITitle:  "inner title",
			content: "inner content",
		},
	}

	//fmt.Println(reflect.Indirect(reflect.ValueOf(m)).FieldByName("ID").CanAddr())
	//fmt.Println(reflect.ValueOf(reflect.Indirect(reflect.ValueOf(m)).FieldByName("ID")).Interface())

	//v := NewItem(m).Fields(`id`, `title`).Resolve()
	//fmt.Printf("%#v",v)
	//vs, _ := json.Marshal(m)
	//fmt.Println(string(vs))
	////
	//fmt.Println(v[`data`].(ResolveMap)[`title`].(string))

	//fmt.Println(v[`data`].(ResolveMap)[`id`].(uint))
	//fmt.Printf("%p", v[`data`].(ResolveMap)[`title`])
}
