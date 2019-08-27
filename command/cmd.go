package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type Cmd interface {
	Use() string
	Run(args []string)
//	Use:   `command`,
//Short: `Firmeve commands`,
//Long:  `Firmeve command collections`,
}



type Console struct {
	root *cobra.Command
}

func (c *Console) Register(cmd Cmd) {
	c.root.AddCommand(&cobra.Command{
		Use:   cmd.Use(),
		//Short: `Firmeve commands`,
		//Long:  `Firmeve command collections`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
		Run: func(cobra *cobra.Command, args []string) {
			cmd.Run(args)
		},
	})
}

func (c *Console) Execute() {
	if err := c.root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var console *Console

func NewCmd() *Console {
	if console != nil {
		return console
	}

	//cobra.OnInitialize()
	console = &Console{
		root: &cobra.Command{
			Use:   `command`,
			Short: `Firmeve commands`,
			Long:  `Firmeve command collections`,
			// Uncomment the following line if your bare application
			// has an action associated with it:
			//	Run: func(cmd *cobra.Command, args []string) { },
		},
	}
	console.root.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return console
}
