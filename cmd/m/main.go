package main

import (
	"fmt"
	http2 "net/http"
	"os"

	"github.com/firmeve/firmeve/http"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/firmeve/firmeve/cmd"
)

type testCmd struct{}

func (t *testCmd) Use() string {
	return "http:serve"
}

func (t *testCmd) Short() string {
	return "Run http server "
}

func (t *testCmd) SetFlags(flags *pflag.FlagSet) {
	flags.IntP("port", "p", 28082, "http port")
}

func (t *testCmd) Handle(cmd *cobra.Command, args []string, flags *pflag.FlagSet) {
	fmt.Println(cmd.Name())
	fmt.Println("My name is main")
}

func (t *testCmd) Register(cmd cmd.ICommand) {

}

type subTestCmd struct{}

func (t *subTestCmd) Use() string {
	return "run"
}

func (t *subTestCmd) Short() string {
	return "Run http server "
}

func (t *subTestCmd) SetFlags(flags *pflag.FlagSet) {
	flags.IntP("port", "p", 28082, "http port")
}

func (t *subTestCmd) Handle(cmd *cobra.Command, args []string, flags *pflag.FlagSet) {
	fmt.Println("My name is main")
}

func main() {
	run := &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			router := http.New()
			router.GET("/abc", func(ctx *http.Context) {
				ctx.String("abc")
				ctx.Next()
			})
			http2.ListenAndServe(cmd.Flag("host").Value.String(), router)
		},
	}
	//run.PersistentFlags().IntP("port", "p", 28188, "http port")
	run.PersistentFlags().StringP("host", "o", "127.0.0.1:28188", "http host")

	http := &cobra.Command{
		Use: "http:serve",
	}
	http.AddCommand(run)

	root := cmd.Root()
	root.AddCommand(http)
	root.SetArgs(os.Args[1:])
	root.Execute()

	//cmd := cmd.New()
	//testCmd := &testCmd{}
	//cmd.Register(testCmd)
	//testCmd.Register()
	//cmd.Root().SetArgs(os.Args[1:])
	//cmd.Execute()
}
