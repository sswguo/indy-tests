package datest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/commonjava/indy-tests/pkg/common"
)

type Report struct {
	ExecutionRoot struct {
		GroupID     string `json:"groupId"`
		ArtifactID  string `json:"artifactId"`
		Version     string `json:"version"`
		OriginalGAV string `json:"originalGAV"`
	} `json:"executionRoot"`

	Modules []struct {
		Gav struct {
			GroupID     string `json:"groupId"`
			ArtifactID  string `json:"artifactId"`
			Version     string `json:"version"`
			OriginalGAV string `json:"originalGAV"`
		} `json:"gav"`
		ManagedPlugins struct {
		} `json:"managedPlugins"`
		ManagedDependencies struct {
			Dependencies map[string]struct {
				GroupID    string `json:"groupId"`
				ArtifactID string `json:"artifactId"`
				Version    string `json:"version"`
			} `json:"dependencies"`
		} `json:"managedDependencies,omitempty"`
	} `json:"modules"`
}

func lookupMetadata(url string) {
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/xml")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Message:") {
		fmt.Print(bodyString)
	}
}

func Run(targetIndy, daGroup string, dataDir string, processNum int) {

	indyHost, validated := common.ValidateTargetIndy(targetIndy)
	if !validated {
		os.Exit(1)
	}

	indyURL := "http://" + indyHost

	routines := processNum

	var urls []string

	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())

		jsonFile, err := os.Open(dataDir + f.Name())

		if err != nil {
			panic(err)
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var report Report

		json.Unmarshal(byteValue, &report)

		fmt.Println("Modules length: ", len(report.Modules))

		for _, module := range report.Modules {
			for _, element := range module.ManagedDependencies.Dependencies {

				groupId := element.GroupID
				artifactId := element.ArtifactID

				fmt.Println("GroupID: ", groupId, " ArtifactId: ", artifactId)

				groupIdPath := strings.ReplaceAll(groupId, ".", "/")

				url := fmt.Sprintf("%s/api/content/maven/group/%s/%s/%s/maven-metadata.xml", indyURL, daGroup, groupIdPath, artifactId)

				urls = append(urls, url)
			}
		}
	}

	fmt.Println("Total requests: ", len(urls), "with routines:", routines)

	concurrentGoroutines := make(chan struct{}, routines)
	var wg sync.WaitGroup

	for i := 0; i < len(urls); i++ {
		concurrentGoroutines <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("Doing", i)
			start := time.Now()
			lookupMetadata(urls[i])
			elapsed := time.Since(start)
			fmt.Println("Finished #", i, " in ", elapsed)
			<-concurrentGoroutines
		}(i)
	}

	wg.Wait()

}
