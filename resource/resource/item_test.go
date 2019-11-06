package resource

import (
	"encoding/json"
	"fmt"
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

	//fmt.Println(reflect.Indirect(reflect.ValueOf(m)).FieldByName("ID").Addr().Interface())
	//fmt.Println(reflect.ValueOf(reflect.Indirect(reflect.ValueOf(m)).FieldByName("ID")).Interface())

	v := NewItem(m).Fields(`id`, `title`).Resolve()

	for k, value := range v[`data`].(ResolveMap) {
		fmt.Println(k,value)
	}

	fmt.Println("==========================")

	//fmt.Printf("%#v",v)
	vs, _ := json.Marshal(v)
	fmt.Println(string(vs))
	////
	//fmt.Printf("%s", v[`data`].(ResolveMap)[`title`])

	//fmt.Println(v[`data`].(ResolveMap)[`id`].(uint))
	//fmt.Printf("%p", v[`data`].(ResolveMap)[`title`])
}
