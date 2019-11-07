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
	return &Command{
		root: &cobra.Command{
			Use:   "firmeve",
			Short: "A generator for Cobra based Applications",
			Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		},
	}
}

func (r *Command) Register(cmd ICmd) {
	r.root.AddCommand(&cobra.Command{
		//Run: cmd.Handle,
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cobraCmd *cobra.Command, args []string) {
			cmd.Handle(r, args)
			fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	})
}

func (r *Command) Root() *cobra.Command {
	return r.root
}

func (r *Command) Execute() {
	err := r.root.Execute()
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
