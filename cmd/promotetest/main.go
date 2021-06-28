package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	exec := &cobra.Command{
		Use:   "indy-promote-test $logUrl",
		Short: "indy-promote-test $logUrl",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	if err := exec.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
