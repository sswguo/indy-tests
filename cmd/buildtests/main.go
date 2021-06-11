package main

import (
	"fmt"
	"os"

	build "commonjava/indy/tests/buildtest"

	"github.com/spf13/cobra"
)

// example: http://orchhost/pnc-rest/v2/builds/97241/logs/build
var logUrl, targetIndy string

func main() {

	exec := &cobra.Command{
		Use:   "indy-build-tests $logUrl",
		Short: "indy-build-tests $logUrl",
		Run: func(cmd *cobra.Command, args []string) {
			if !validate(args) {
				cmd.Help()
				os.Exit(1)
			}
			logUrl = args[0]
			run(logUrl)
		},
	}

	exec.Flags().StringVarP(&targetIndy, "targetIndy", "t", "", "The target indy server to do the testing. If not specified, will get from env variables 'INDY_TARGET'")

	if err := exec.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validate(args []string) bool {
	if len(args) <= 0 {
		fmt.Printf("logUrl is not specified!\n\n")
		return false
	}
	if build.IsEmptyString(args[0]) {
		fmt.Printf("logUrl cannot be empty!\n\n")
		return false
	}
	if build.IsEmptyString(targetIndy) {
		targetIndy = os.Getenv("INDY_TARGET")
		if build.IsEmptyString(targetIndy) {
			fmt.Printf("The target indy server can not be empty!\n\n")
			return false
		}
	}
	return true
}

func run(logUrl string) {
	log, err := build.GetRespAsPlaintext(logUrl)
	if err != nil {
		httpErr := err.(build.HTTPError)
		fmt.Printf("Request failed! Log url: %s, response status: %d, error message: %s\n", logUrl, httpErr.StatusCode, httpErr.Message)
		os.Exit(1)
	}
	result, err := build.ParseLog(log)
	if err != nil {
		fmt.Printf("Log parse failed! Log url: %s, error message: %s\n", logUrl, err.Error())
		os.Exit(1)
	}
	if err == nil {
		downloads := result["downloads"]
		if downloads != nil {
			fmt.Println("Start showing downloads: ==================\n")
			for _, d := range downloads {
				fmt.Println(d)
			}
			fmt.Println("\nFinish showing downloads: ==================\n")
		}
		uploads := result["uploads"]
		if uploads != nil {
			fmt.Println("Start showing uploads: ==================\n")
			for _, u := range uploads {
				fmt.Println(u)
			}
			fmt.Println("\nFinish showing uploads: ==================\n")
		}
	}
}
