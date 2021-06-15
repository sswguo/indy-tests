package buildtest

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecorate(t *testing.T) {
	Convey("Decorate should work", t, func() {
		Convey("No duplication", func() {
			testDownloads := []string{"http://testdownload/maven/for/test/org/jboss/eap/wildfly-clustering-web-extension/7.3.8.GA-redhat-00001/wildfly-clustering-web-extension-7.3.8.GA-redhat-00001.jar",
				"http://testdownload/maven/for/test/org/jboss/eap/wildfly-configadmin/7.3.8.GA-redhat-00001/wildfly-configadmin-7.3.8.GA-redhat-00001.pom"}
			finalDowns := decorateChecksums(testDownloads)
			So(len(finalDowns), ShouldEqual, 8)
		})
		Convey("Should avoid duplication", func() {
			testDownloads := []string{"http://testdownload/maven/for/test/org/jboss/eap/wildfly-clustering-web-extension/7.3.8.GA-redhat-00001/wildfly-clustering-web-extension-7.3.8.GA-redhat-00001.jar",
				"http://testdownload/maven/for/test/org/jboss/eap/wildfly-configadmin/7.3.8.GA-redhat-00001/wildfly-configadmin-7.3.8.GA-redhat-00001.pom",
				"http://testdownload/maven/for/test/org/jboss/eap/wildfly-configadmin/7.3.8.GA-redhat-00001/wildfly-configadmin-7.3.8.GA-redhat-00001.pom.md5",
				"http://testdownload/maven/for/test/org/jboss/eap/wildfly-configadmin/7.3.8.GA-redhat-00001/wildfly-configadmin-7.3.8.GA-redhat-00001.pom.sha1"}
			finalDowns := decorateChecksums(testDownloads)
			So(len(finalDowns), ShouldEqual, 8)
		})

	})
}

func TestReplaceTarget(t *testing.T) {
	Convey("Replacing should work", t, func() {
		testDownloads := []string{"http://testdownload/maven/for/test/org/jboss/eap/wildfly-clustering-web-extension/7.3.8.GA-redhat-00001/wildfly-clustering-web-extension-7.3.8.GA-redhat-00001.jar",
			"http://testdownload/maven/for/test/org/jboss/eap/wildfly-configadmin/7.3.8.GA-redhat-00001/wildfly-configadmin-7.3.8.GA-redhat-00001.pom"}
		replaced := replaceTarget(testDownloads, "", "replacedtestdownload")
		So(replaced[0], ShouldEqual, "http://replacedtestdownload/maven/for/test/org/jboss/eap/wildfly-clustering-web-extension/7.3.8.GA-redhat-00001/wildfly-clustering-web-extension-7.3.8.GA-redhat-00001.jar")
		So(replaced[1], ShouldEqual, "http://replacedtestdownload/maven/for/test/org/jboss/eap/wildfly-configadmin/7.3.8.GA-redhat-00001/wildfly-configadmin-7.3.8.GA-redhat-00001.pom")
	})
}
