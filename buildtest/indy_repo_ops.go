package buildtest

import (
	"fmt"
	"io"
	"strings"
)

const (
	DEFAULT_SHARED_GROUP = "builds-untested+shared-imports+public"
	DEFAULT_MVN_CENTRAL  = "central"
	DEFAULT_NPM_CENTRAL  = "npmjs"
	TYPE_MVN             = "maven"
	TYPE_NPM             = "npm"
)

type BuildMetadata struct {
	buildType     string
	centralName   string
	sharedGrpName string
}

func decideMeta(buildType string) *BuildMetadata {
	if buildType == TYPE_MVN {
		return &BuildMetadata{
			buildType:     buildType,
			centralName:   DEFAULT_MVN_CENTRAL,
			sharedGrpName: DEFAULT_SHARED_GROUP,
		}
	} else if buildType == TYPE_NPM {
		return &BuildMetadata{
			buildType:   buildType,
			centralName: DEFAULT_NPM_CENTRAL,
		}
	}
	return nil
}

func prepareIndyRepos(indyURL, buildName string, buildMeta BuildMetadata) bool {
	if !prepareIndyHosted(indyURL, buildMeta.buildType, buildName) {
		return true
	}
	return prepareIndyGroup(indyURL, buildName, buildMeta)
}

func prepareIndyHosted(indyURL, buildType, buildName string) bool {
	hostedVars := IndyHostedVars{
		Name: buildName,
		Type: buildType,
	}

	URL := fmt.Sprintf("%s/api/admin/stores/%s/hosted/%s", indyURL, buildType, buildName)

	hosted := IndyHostedTemplate(&hostedVars)
	fmt.Printf("Start creating hosted repo %s\n", buildName)
	result := putRequest(URL, strings.NewReader(hosted))
	if result {
		fmt.Printf("Hosted repo %s created successfully, check %s for details\n", buildName, URL)
	} else {
		fmt.Printf("Hosted repo %s creation fail, no following operations\n", buildName)
	}
	return result
}

func prepareIndyGroup(indyURL, buildName string, buildMeta BuildMetadata) bool {
	buildType, sharedGrpName := buildMeta.buildType, buildMeta.sharedGrpName
	groupVars := IndyGroupVars{
		Name:         buildName,
		Type:         buildMeta.buildType,
		Constituents: []string{fmt.Sprintf("%s:hosted:%s", buildType, buildName), fmt.Sprintf("%s:group:%s", buildType, sharedGrpName)},
	}
	group := IndyGroupTemplate(&groupVars)

	URL := fmt.Sprintf("%s/api/admin/stores/%s/group/%s", indyURL, buildType, buildName)

	fmt.Printf("Start creating group repo %s\n", buildName)
	result := putRequest(URL, strings.NewReader(group))
	if result {
		fmt.Printf("Group repo %s created successfully, check %s for details\n", buildName, URL)
	} else {
		fmt.Printf("Group repo %s created failed, no following operations\n", buildName)
	}
	return result
}

func destroyIndyRepos(indyURL, buildType, buildName string) {
	destroyIndyGroup(indyURL, buildType, buildName)
	// destroyHosted(indyURL, buildName)
}

func destroyIndyHosted(indyURL, buildType, buildName string) {
	URL := fmt.Sprintf("%s/api/admin/stores/%s/hosted/%s", indyURL, buildType, buildName)
	fmt.Printf("Start deleting hosted repo %s\n", buildName)
	result := delRequest(URL)
	if result {
		fmt.Printf("Hosted repo %s deleted successfully\n", buildName)
	}
}

func destroyIndyGroup(indyURL, buildType, buildName string) {
	URL := fmt.Sprintf("%s/api/admin/stores/%s/group/%s", indyURL, buildType, buildName)
	fmt.Printf("Start deleting group repo %s\n", buildName)
	result := delRequest(URL)
	if result {
		fmt.Printf("Group repo %s deleted successfully\n", buildName)
	}
}
func getRequest(url string) (string, bool) {
	content, _, succeeded := HTTPRequest(url, MethodGet, nil, true, nil, nil, "", false)
	return content, succeeded
}

func postRequest(url string, data io.Reader) (string, bool) {
	content, _, succeeded := HTTPRequest(url, MethodPost, nil, true, data, nil, "", false)
	return content, succeeded
}

func putRequest(url string, data io.Reader) bool {
	_, _, succeeded := HTTPRequest(url, MethodPut, nil, false, data, nil, "", false)
	return succeeded
}

func delRequest(url string) bool {
	_, _, succeeded := HTTPRequest(url, MethodDelete, nil, false, nil, nil, "", false)
	return succeeded
}
