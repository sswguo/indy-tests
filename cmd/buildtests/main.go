/*
 *  Copyright (C) 2011-2021 Red Hat, Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *          http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	build "github.com/commonjava/indy-tests/buildtest"
	"github.com/commonjava/indy-tests/common"

	"github.com/spf13/cobra"
)

// example: http://orchhost/pnc-rest/v2/builds/97241/logs/build
var logUrl, targetIndy, repoReplPattern, buildType string
var processNum int

const DEFAULT_PROCESS_NUM = 1
const DEFAULT_REPO_REPL_PATTERN = ""
const DEFAULT_BUILD_TYPE = "maven"

func main() {

	exec := &cobra.Command{
		Use:   "indy-build-tests $logUrl",
		Short: "indy-build-tests $logUrl",
		Run: func(cmd *cobra.Command, args []string) {
			if !validate(args) {
				cmd.Help()
				os.Exit(1)
			}
			// here will use env variables if they are specified for some flags
			checkEnvVars()
			logUrl = args[0]
			build.Run(logUrl, "", targetIndy, buildType, processNum)
		},
	}

	exec.Flags().StringVarP(&targetIndy, "targetIndy", "t", "", "The target indy server to do the testing. If not specified, will get from env variables 'INDY_TARGET'.")
	exec.Flags().StringVarP(&buildType, "buildType", "b", DEFAULT_BUILD_TYPE, "The type of the build, should be 'maven' or 'npm'. Default is 'maven'.")
	exec.Flags().IntVarP(&processNum, "processNum", "p", DEFAULT_PROCESS_NUM, "The number of processes to download and upload files in parralel.")

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
	if common.IsEmptyString(args[0]) {
		fmt.Printf("logUrl cannot be empty!\n\n")
		return false
	}
	if common.IsEmptyString(targetIndy) {
		targetIndy = os.Getenv("INDY_TARGET")
		if common.IsEmptyString(targetIndy) {
			fmt.Printf("The target indy server can not be empty!\n\n")
			return false
		}
	}

	return true
}

func checkEnvVars() {
	envBuildType := os.Getenv("INDY_BUILD_TYPE")
	if envBuildType == "maven" || envBuildType == "npm" {
		buildType = envBuildType
	}
	envProcNum := os.Getenv("BUILD_PROC_NUM")
	if num, err := strconv.Atoi(envProcNum); err == nil {
		processNum = num
	}
}
