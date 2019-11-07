package cmd

import (
	"fmt"

	"github.com/firmeve/firmeve"
	"github.com/spf13/cobra"
)

type ICmd interface {
	Handle(cmd *Command, args []string)
}

type Command struct {
	root    *cobra.Command
	Firmeve *firmeve.Firmeve
}

func New() *Command {
	root := &cobra.Command{
		Use:   "firmeve",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
	//root.DisableSuggestions = true
	root.SetVersionTemplate("main version 11111")
	root.PersistentFlags().StringP("version", "v", "", "author name for copyright attribution")
	return &Command{
		root: root,
	}
}

func (r *Command) Register(cmd ICmd) {
	c := &cobra.Command{
		//Run: cmd.Handle,
		Use:   "version",
		Short: "Current version 1.0",
		//Long:  `All software has versions. This is Hugo's`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
		},
		Run: func(cobraCmd *cobra.Command, args []string) {
			//cmd.Handle(r, args)
			fmt.Println("1111")
			//fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}
	c.PersistentFlags().StringP("cccccccc", "c", "", "SSSS")
	//c.SetHelpTemplate("ssssssssssssssssssssssssss111\nabc\n")
	c.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println("GGGGGG")
		return nil
	})
	c.SetUsageTemplate("111111")
	//c.SetHelpCommand(c)
	r.root.AddCommand(c)
}

func (r *Command) Root() *cobra.Command {
	return r.root
}

func (r *Command) Execute() {
	err := r.root.Execute()
	if err != nil {
		//panic(err)
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
