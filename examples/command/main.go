package main

import (
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/cobra"
)

type TestCommand struct {
}

func (t TestCommand) CobraCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "testing"
	cmd.Short = "Testing a cmd"
	return cmd
}

func (t TestCommand) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	fmt.Println("run")
}

func main() {
	firmeve.RunDefault(firmeve.WithConfigPath(`../config.yaml`), firmeve.WithCommands(
		[]contract.Command{
			new(TestCommand),
		},
	))
}
