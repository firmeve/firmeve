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
	Content string `resource:"_content"`
	SubMock *innerMock `resource:"inner_mock"`
}

func TestNewItem(t *testing.T) {
	m := &mock{
		ID:10,
		Title:"main title",
		Content:"content",
		SubMock:&innerMock{
			IId:11,
			ITitle:"inner title",
			content:"inner content",
		},
	}
	v:=NewItem(m).Fields(`id`,`title`).Resolve()
	vs,_ := json.Marshal(v)
	fmt.Println(string(vs))

	//fmt.Println(v[`data`].(ResolveMap)[`id`].(uint))
	fmt.Printf("%v",len(v[`data`].(ResolveMap)[`title`].(string)))
}