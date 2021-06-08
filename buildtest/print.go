/*
 *  Copyright (C) 2011-2020 Red Hat, Inc.
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

package buildtest

import (
	"fmt"

	"github.com/jedib0t/go-pretty/text"
)

const textWidth = 25

func Printlnf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func PrintProps(name string, value string) {
	// trace.Printf("key:%s : value: %s", name, value)
	if !IsEmptyString(name) && !IsEmptyString(value) {
		fmt.Printf("%s:   %s\n", text.AlignLeft.Apply(name, textWidth), text.AlignLeft.Apply(value, textWidth))
	}
}

func PrintVerbose(message string, verbose bool) {
	if verbose {
		fmt.Print(message)
	}
}
