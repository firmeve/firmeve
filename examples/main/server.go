package main

import (
	"fmt"
	"os"

	"github.com/firmeve/firmeve/converter/resource"

	"github.com/firmeve/firmeve"

	"github.com/firmeve/firmeve/cmd"
	"github.com/firmeve/firmeve/converter/transform"
	"github.com/firmeve/firmeve/http"
	_ "github.com/takama/daemon"
)

type Something struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func (s *Something) SetId() {
	s.Id += 1
}

type Original struct {
	Id    int
	Title string
	Name  string
}

type AnyTransformer struct {
	transform.BaseTransformer
}

func (a *AnyTransformer) IdField() int {
	return a.Resource().(*Original).Id * 10
}

func main() {
	//x := 1
	//y := 2
	//fmt.Println(reflect.TypeOf(&x).Kind())
	//fmt.Println(reflect.TypeOf(&x) == reflect.TypeOf(&y))
	//fmt.Println("#########################")
	firmeve.Instance().Bind("something", func() *Something {
		return firmeve.Instance().Make(new(Something)).(*Something)
	})
	router := http.New()
	router.GET(`/item`, func(c *http.Context) {
		c.Item(&Original{
			10,
			"Title",
			"Simon",
		}, new(AnyTransformer))
	})
	router.GET(`/collection`, func(c *http.Context) {
		z := []*resource.Item{
			resource.NewItem(transform.New(&Original{
				10, "title", "simon",
			}, new(AnyTransformer))),
			resource.NewItem(transform.New(&Original{
				11, "title", "simon",
			}, new(AnyTransformer))),
			resource.NewItem(transform.New(&Original{
				12, "title", "simon",
			}, new(AnyTransformer))),
			resource.NewItem(transform.New(&Original{
				13, "title", "simon",
			}, new(AnyTransformer))),
		}
		//z := []*resource.Item{
		//	resource.NewItem(transform.New(&Original{
		//{
		//	10,
		//	"Title",
		//	"Simon",
		//}, new(AnyTransformer)))()
		//}
		c.Collection(z, new(AnyTransformer))
	})
	router.GET("/abc", func(ctx *http.Context) {
		//time.Sleep(time.Second * 10)
		//ctx.String("abc")
		//s := ctx.Firmeve.Get("something").(*Something)
		//s.SetId()
		//fmt.Printf("%#v\n", s)
		//ctx.Json(s)
		s := ctx.Firmeve.Make(new(Something)).(*Something)
		s.SetId()
		fmt.Printf("%#v\n", s)
		ctx.Json(s)
		//if !ctx.Firmeve.Has("something") {
		//	ctx.Firmeve.Bind("something", ctx.Firmeve.Make(new(Something)))
		//	ctx.Json(Something{
		//		Id:    0,
		//		Title: "",
		//	})
		//} else {
		//	s := ctx.Firmeve.Get("something").(*Something)
		//	s.SetId()
		//	fmt.Printf("%#v\n", s)
		//	ctx.Json(s)
		//}

		ctx.Next()
	}).After(func(ctx *http.Context) {
		//ctx.Json(Something{
		//	Id:    10,
		//	Title: "title",
		//})
		ctx.Next()
	})
	//
	root := cmd.Root()
	root.AddCommand(http.NewCmd(router).Cmd())
	root.SetArgs(os.Args[1:])
	root.Execute()

	//cmd := cmd.New()
	//testCmd := &testCmd{}
	//cmd.Register(testCmd)
	//testCmd.Register()
	//cmd.Root().SetArgs(os.Args[1:])
	//cmd.Execute()
}
