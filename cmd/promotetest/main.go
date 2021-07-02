package main

import (
	"fmt"
	"os"

	"github.com/commonjava/indy-tests/promotetest"
	"github.com/spf13/cobra"
)

var targetIndy, trackingId, promoteTarget string

func main() {

	exec := &cobra.Command{
		Use:   "indy-promote-test $targetIndy $trackingId $promoteTarget",
		Short: "indy-promote-test $targetIndy $trackingId $promoteTarget",
		Run: func(cmd *cobra.Command, args []string) {
			if !validate(args) {
				cmd.Help()
				os.Exit(1)
			}

			promotetest.Run(args[0], args[1], args[2])
		},
	}

	if err := exec.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validate(args []string) bool {
	if len(args) <= 2 {
		fmt.Printf("there are at least 3 non-empty arguments: targetIndy, trackingId, promoteTarget!\n\n")
		return false
	}
	return true
}
