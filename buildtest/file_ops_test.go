package buildtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileExists(t *testing.T) {
	Convey("TestFileExists", t, func() {
		Convey("File should exists", func() {
			So(FileOrDirExists("/usr/bin/bash"), ShouldBeTrue)
		})
		Convey("File should not exists", func() {
			So(FileOrDirExists("/kljsdflksdjf"), ShouldBeFalse)
		})
	})

}

func TestStoreFile(t *testing.T) {
	Convey("TestStoreFile", t, func() {
		fileName := fmt.Sprintf("/tmp/%d", nowInMillis())
		fileContent := "This is a test."
		StoreFile(fileName, fileContent, true)
		Convey("Stored file should exist", func() {
			So(FileOrDirExists(fileName), ShouldBeTrue)
		})
		f, _ := os.Open(fileName)
		defer f.Close()
		actual, _ := ioutil.ReadAll(f)
		Convey("Stored file content should be correct", func() {
			So(string(actual), ShouldEqual, fileContent)
		})
	})

}
func nowInMillis() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
