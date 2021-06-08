package main

import (
	build "commonjava/indy/tests/buildtest"
	"fmt"
)

func main() {
	log, _ := build.GetRespAsPlaintext("http://orchhost/pnc-rest/v2/builds/97241/logs/build")
	result, err := build.ParseLog(log)
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
