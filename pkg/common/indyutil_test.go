package common

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStoreKeyToPath(t *testing.T) {
	Convey("TestFileExists", t, func() {
		So(StoreKeyToPath("maven:hosted:shared-imports"), ShouldEqual, "maven/hosted/shared-imports")
	})
}
