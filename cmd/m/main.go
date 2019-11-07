package main

import (
	"github.com/firmeve/firmeve/cmd"
	http2 "github.com/firmeve/firmeve/http"
	"net/http"
)

type testCmd struct {
}

func (t *testCmd) Handle(cmd *cmd.Command, args []string) {
	cmd.Root().Short = "sssss"
	//fmt.Println("My name is main")
	router := http2.New()
	router.GET("/gets/1", func(ctx *http2.Context) {
		ctx.Write([]byte("Body"))
		ctx.Next()
	})
	err := http.ListenAndServe("0.0.0.0:28082", router)
	if err != nil {
		panic(err)
	}
}

func main() {
	cmd := cmd.New()
	cmd.Register(&testCmd{})
	cmd.Execute()
}
