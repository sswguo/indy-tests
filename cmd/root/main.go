package main

import (
	"fmt"
	"os"

	"github.com/commonjava/indy-tests/cmd/buildtest"
	"github.com/commonjava/indy-tests/cmd/promotetest"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "indy-test",
		Short: "indy-test is a tool to do indy integration test against runnable indy server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(buildtest.NewBuildTestCmd())
	rootCmd.AddCommand(promotetest.NewPromoteTestCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
