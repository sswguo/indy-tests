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

package integrationtest

import (
	"fmt"
	"os"

	"github.com/commonjava/indy-tests/pkg/common"
)

const TARGET_DIR = "target"

/*
 * When we start the integration test, it will (in order):
 *
 * a. Load the corresponding dataset
 * b. Retrieve the metadata files in da.json to make sure they can be downloaded successfully
 * c. Create a temp group which is same to PNC temp build group (with a hosted repo, namely A)
 * d. Download the files in tracking "downloads" section from this temp group
 * e. Download the files in tracking "uploads" section, and rename all the files (jar, pom, and so on)
 *    with a new version suffix and upload them again to the hosted repo A
 * f. Retrieve the metadata files that will be updated by below promotion (we need to foresee those metadata files)
 * g. Promote the files in hosted repo A to hosted repo pnc-builds.
 *    Because we rename the artifacts and pom files, we can run a test multiple times without conflicting each other.
 * h. Retrieve the metadata files from step #f again, check if the new version is available
 * i. (optional) Delete the temp group and the hosted repo A. This is not mandatory because we use renamed versions.
 *    Leaving them there won't affect the following tests.
 */
func Run(indyBaseUrl, datasetRepoUrl, buildId string) {
	//Create target folder (to store downloaded files), e.g, 'target'
	err := os.MkdirAll(TARGET_DIR, 0755)
	check(err)

	//Clone dataset repo
	datasetRepoDir := funcA_CloneRepo(datasetRepoUrl)
	fmt.Printf("Clone SUCCESS, dir: %s\n", datasetRepoDir)

	//TODO: Retrieve the metadata files in da.json
	funcB()

	funcC()

	funcD()

	funcE()

	funcF()

	funcG()

	funcH()

	funcI()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func funcA_CloneRepo(datasetRepoUrl string) string {
	return common.DownloadRepo(datasetRepoUrl)
}

func funcB() {

}

func funcC() {

}

func funcD() {

}

func funcE() {

}

func funcF() {

}

func funcG() {

}

func funcH() {

}

func funcI() {

}
