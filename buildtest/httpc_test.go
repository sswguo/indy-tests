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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHost(t *testing.T) {
	Convey("GetHost should return correct hostname", t, func() {
		Convey("https://www.google.com host is www.google.com", func() {
			So(GetHost("https://www.google.com"), ShouldEqual, "www.google.com")
		})
		Convey("http://www.test.com host is www.test.com", func() {
			So(GetHost("http://www.test.com"), ShouldEqual, "www.test.com")
		})
	})
}

func TestGetPort(t *testing.T) {
	Convey("GetPort should return correct port", t, func() {
		Convey("https://www.google.com port is empty", func() {
			So(GetPort("https://www.google.com"), ShouldEqual, "")
		})
		Convey("http://www.test.com:8080 port is 8080", func() {
			So(GetPort("http://www.test.com:8080"), ShouldEqual, "8080")
		})
	})
}
