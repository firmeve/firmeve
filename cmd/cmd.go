package cmd

import (
	reflect2 "reflect"

	"github.com/firmeve/firmeve/support/reflect"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

type ICommand interface {
	Use() string
	SetFlags(flags *pflag.FlagSet)
	Handle(cmd *cobra.Command, args []string, flags *pflag.FlagSet)
	//New(cmd *Command) *cobra.Command

	//Handle(cmd *Command, args []string, flags *flag.FlagSet)
}

func Root() *cobra.Command {
	return &cobra.Command{
		Use:   "firmeve",
		Short: "Firmeve command",
		Long:  "Firmeve",
	}
}

type Command struct {
	root *cobra.Command
}

func New() *Command {
	cmd := &Command{}
	cmd.Register(cmd)
	return cmd
}

func (c *Command) Use() string {
	return `firmeve [options] command`
}

func (c *Command) Short() string {
	return `firmeve command`
}

func (c *Command) Long() string {
	return `firmeve long command`
}

func (c *Command) SetFlags(flags *pflag.FlagSet) {
	flags.StringP("version", "v", "", "author name for copyright attribution")
}

func (c *Command) Handle(cmd *cobra.Command, args []string, flags *pflag.FlagSet) {
	//fmt.Println(cmd.Name(), args, flags)
	//	cmd.Printf("%s", `
	//Usage:
	//  firmeve [command]
	//
	//Available Commands:
	//  help        Help about any command
	//  version
	//
	//Flags:
	//  -h, --help   help for firmeve
	//
	//Use "firmeve [command] --help" for more information about a command.
	//`)
}

func (c *Command) Register(command ICommand) {
	cmd := new(cobra.Command)
	cmd.Use = command.Use()
	reflectType := reflect2.TypeOf(command)
	reflectValue := reflect2.ValueOf(command)
	if reflect.MethodExists(reflectType, "Short") {
		cmd.Short = reflect.CallMethodValue(reflectValue, "Short")[0].(string)
	}
	if reflect.MethodExists(reflectType, "Long") {
		cmd.Long = reflect.CallMethodValue(reflectValue, "Long")[0].(string)
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		command.Handle(cmd, args, cmd.Flags())
	}
	command.SetFlags(cmd.PersistentFlags())
	if c.root == nil {
		c.root = cmd
		c.root.Run = nil
	} else {

		c.root.AddCommand(cmd)
	}

	//cmd.Short
	//c := &cobra.Command{
	//	//Run: cmd.Handle,
	//	Use:   "version",
	//	Short: "Current version 1.0",
	//	//Long:  `All software has versions. This is Hugo's`,
	//	//PersistentPreRun: func(cmd *cobra.Command, args []string) {
	//	//	fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
	//	//},
	//	//PreRun: func(cmd *cobra.Command, args []string) {
	//	//	fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
	//	//},
	//	//PostRun: func(cmd *cobra.Command, args []string) {
	//	//	fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
	//	//},
	//	//PersistentPostRun: func(cmd *cobra.Command, args []string) {
	//	//	fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
	//	//},
	//	Run: func(cobraCmd *cobra.Command, args []string) {
	//
	//		//cmd.Handle(r, args)
	//		fmt.Println(cobraCmd.Use)
	//		fmt.Println(cobraCmd.Flag("ccc").Value.String())
	//		//fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	//	},
}

//c.PersistentFlags().StringP("ccc", "c", "", "SSSS")
//c.SetHelpTemplate("ssssssssssssssssssssssssss111\nabc\n")
//c.SetUsageFunc(func(command *cobra.Command) error {
//	fmt.Println("GGGGGG")
//	return nil
//})
//c.SetUsageTemplate("111111")
//c.SetHelpCommand(c)
//r.root.RemoveCommand(c)
//r.Commands().
//}

func (c *Command) Root() *cobra.Command {
	return c.root
}

func (c *Command) Execute() {
	err := c.root.Execute()
	if err != nil {
		panic(err)
	}
}

//
//var rootCmd = &cobra.Command{
//	Use:   "hugo",
//	Short: "Hugo is a very fast static site generator",
//	Long: `A Fast and Flexible Static Site Generator built with
//                love by spf13 and friends in Go.
//                Complete documentation is available at http://hugo.spf13.com`,
//	Run: func(cmd *cobra.Command, args []string) {
//		// Do Stuff Here
//	},
//}
//
//func Execute() {
//	if err := rootCmd.Execute(); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//}
