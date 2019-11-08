package main

import (
	"os"
	"time"

	"github.com/firmeve/firmeve/http"

	"github.com/firmeve/firmeve/cmd"
	_ "github.com/takama/daemon"
)

func main() {

	router := http.New()
	router.GET("/abc", func(ctx *http.Context) {
		time.Sleep(time.Second * 10)
		ctx.String("abc")
		ctx.Next()
	})

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
