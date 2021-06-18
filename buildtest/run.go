package buildtest

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const TMP_DOWNLOAD_DIR = "/tmp/download"

func Run(logUrl, replacement, targetIndy, buildType string, processNum int) {
	indyHost, validated := validateTargetIndy(targetIndy)
	if !validated {
		os.Exit(1)
	}

	newBuildName := generateRandomBuildName()

	//TODO: enable this one when in a working testing indy env
	// buildMeta := decideMeta(buildType)
	// if !prepareIndyRepos("http://"+indyHost, newBuildName, *buildMeta) {
	// 	os.Exit(1)
	// }

	log, err := GetRespAsPlaintext(logUrl)
	if err != nil {
		httpErr := err.(HTTPError)
		fmt.Printf("Request failed! Log url: %s, response status: %d, error message: %s\n", logUrl, httpErr.StatusCode, httpErr.Message)
		os.Exit(1)
	}
	result, err := ParseLog(log)
	if err != nil {
		fmt.Printf("Log parse failed! Log url: %s, error message: %s\n", logUrl, err.Error())
		os.Exit(1)
	}

	if !FileOrDirExists(TMP_DOWNLOAD_DIR) {
		os.Mkdir(TMP_DOWNLOAD_DIR, os.FileMode(0755))
	}
	if !FileOrDirExists(TMP_DOWNLOAD_DIR) {
		fmt.Printf("Error: cannot create directory %s for file downloading.\n", TMP_DOWNLOAD_DIR)
		os.Exit(1)
	}
	if err == nil {
		downloads := replaceTarget(decorateChecksums(result["downloads"]), "", indyHost, newBuildName)
		result["downloads"] = nil
		if downloads != nil {
			if processNum > 1 {
				//TODO: implement concurrent download here with processNum
			} else {
				for _, url := range downloads {
					fileLoc := path.Join(TMP_DOWNLOAD_DIR, path.Base(url))
					DownloadFile(url, fileLoc)
				}
			}
		}
		//TODO: handle uploads here
		// uploads := replaceTarget(result["uploads"], "", targetIndy)
		// result["uploads"] = nil
		// if uploads != nil {
		// 	fmt.Print("Start showing uploads: ==================\n\n")
		// 	for _, u := range uploads {
		// 		fmt.Println(u)
		// 	}
		// 	fmt.Print("\nFinish showing uploads: ==================\n\n")
		// }
	}
}

func validateTargetIndy(targetIndy string) (string, bool) {
	indyHost := targetIndy
	if strings.HasPrefix(targetIndy, "http") {
		host := GetHost(targetIndy)
		port := GetPort(targetIndy)
		if IsEmptyString(port) || port == "80" {
			indyHost = host
		} else {
			indyHost = fmt.Sprintf("%s:%s", host, port)
		}
	} else {
		indyHost = strings.Split(targetIndy, "/")[0]
	}

	indyTest := ""
	var indyAPIBase string
	if strings.HasPrefix(indyHost, "http") {
		indyAPIBase = indyHost + "/api"
	} else {
		indyAPIBase = "http://" + indyHost + "/api"
	}

	Printlnf("Start testing target indy server %s", indyHost)
	_, err := url.ParseRequestURI(indyAPIBase)
	if err == nil {
		testPath := "/admin/stores/maven/remote/central"
		indyTest = indyAPIBase + testPath
		_, err = url.ParseRequestURI(indyTest)
		if err != nil {
			Printlnf("Error: not a valid indy server: %s because %s does not exist", targetIndy, testPath)
			return "", false
		}
	} else {
		Printlnf("Error: not a valid indy server: %s", targetIndy)
		return "", false
	}
	resp, err2 := http.Get(indyTest)
	if err2 != nil {
		Printlnf("Error: %s is not a valid indy server. Cause: %s", targetIndy, err2)
		return "", false
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		Printlnf("Error: %s returned bad status. Cause: %s", targetIndy, resp.Status)
		return "", false
	}
	resp.Body.Close()
	return indyHost, true
}

func decorateChecksums(downloads []string) []string {
	downSet := make(map[string]bool)
	for _, artifact := range downloads {
		downSet[artifact] = true
		if strings.HasSuffix(artifact, ".md5") || strings.HasSuffix(artifact, ".sha1") || strings.HasSuffix(artifact, ".sha256") {
			continue
		}
		downSet[artifact+".md5"] = true
		downSet[artifact+".sha1"] = true
		downSet[artifact+".sha256"] = true
	}
	finalDownloads := []string{}
	for artifact := range downSet {
		finalDownloads = append(finalDownloads, artifact)
	}
	return finalDownloads
}

func replaceTarget(artifacts []string, oldIndyHost, targetIndyHost, buildName string) []string {
	results := []string{}
	for _, a := range artifacts {
		repl := oldIndyHost
		if IsEmptyString(repl) {
			repl = a[strings.Index(a, "//")+2:]
			repl = repl[:strings.Index(repl, "/")]
		}
		final := strings.ReplaceAll(a, repl, targetIndyHost)
		buildPat := regexp.MustCompile(`https{0,1}:\/\/.+\/(build-\d+)\/.*`)
		buildPat.FindAllStringSubmatch(final, 0)
		matches := buildPat.FindAllStringSubmatch(final, -1)
		if matches != nil {
			for i := range matches {
				get := matches[i][1]
				if strings.HasPrefix(get, "build-") {
					final = strings.ReplaceAll(final, get, buildName)
					break
				}
			}
		}

		results = append(results, final)
	}
	return results
}

func generateRandomBuildName() string {
	buildPrefix := "build-test-"
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 99999
	return fmt.Sprintf(buildPrefix+"%v", rand.Intn(max-min)+min)
}
