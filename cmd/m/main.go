package main

import (
	"fmt"
	"os"

	"github.com/firmeve/firmeve"

	"github.com/firmeve/firmeve/cmd"
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
