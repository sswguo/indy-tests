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
const TMP_UPLOAD_DIR = "/tmp/upload"

func Run(logUrl, replacement, targetIndy, buildType string, processNum int) {
	indyHost, validated := validateTargetIndy(targetIndy)
	if !validated {
		os.Exit(1)
	}

	newBuildName := generateRandomBuildName()

	// Prepare the indy repos for the whole testing
	buildMeta := decideMeta(buildType)
	if !prepareIndyRepos("http://"+indyHost, newBuildName, *buildMeta) {
		os.Exit(1)
	}

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

	prepareCacheDirectories()

	if err == nil {
		downloads := replaceTargets(decorateChecksums(result["downloads"]), "", indyHost, newBuildName)
		result["downloads"] = nil // save memory
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
		// uploads := replaceTargets(result["uploads"], "", indyHost, newBuildName)
		// result["uploads"] = nil // save memory
		if result["uploads"] != nil {
			if processNum > 1 {
				//TODO: implement concurrent upload here with processNum
			} else {
				for _, url := range result["uploads"] {
					cacheFile := path.Join(TMP_UPLOAD_DIR, path.Base(url))
					downloadArtifact := replaceHost(url, "", indyHost)
					downloaded := DownloadUploadFileForCache(downloadArtifact, cacheFile)
					if downloaded {
						replacedUrl := replaceBuildName(downloadArtifact, newBuildName)
						UploadFile(replacedUrl, cacheFile)
					}
				}
			}
		}
	}
}

func prepareCacheDirectories() {
	if !fileOrDirExists(TMP_DOWNLOAD_DIR) {
		os.Mkdir(TMP_DOWNLOAD_DIR, os.FileMode(0755))
	}
	if !fileOrDirExists(TMP_DOWNLOAD_DIR) {
		fmt.Printf("Error: cannot create directory %s for file downloading.\n", TMP_DOWNLOAD_DIR)
		os.Exit(1)
	}
	if !fileOrDirExists(TMP_UPLOAD_DIR) {
		os.Mkdir(TMP_UPLOAD_DIR, os.FileMode(0755))
	}
	if !fileOrDirExists(TMP_UPLOAD_DIR) {
		fmt.Printf("Error: cannot create directory %s for caching uploading files.\n", TMP_UPLOAD_DIR)
		os.Exit(1)
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

	fmt.Printf("Start testing target indy server %s\n", indyHost)
	_, err := url.ParseRequestURI(indyAPIBase)
	if err == nil {
		testPath := "/admin/stores/maven/remote/central"
		indyTest = indyAPIBase + testPath
		_, err = url.ParseRequestURI(indyTest)
		if err != nil {
			fmt.Printf("Error: not a valid indy server: %s because %s does not exist\n", targetIndy, testPath)
			return "", false
		}
	} else {
		fmt.Printf("Error: not a valid indy server: %s\n", targetIndy)
		return "", false
	}
	resp, err2 := http.Get(indyTest)
	if err2 != nil {
		fmt.Printf("Error: %s is not a valid indy server. Cause: %s\n", targetIndy, err2)
		return "", false
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Printf("Error: %s returned bad status. Cause: %s\n", targetIndy, resp.Status)
		return "", false
	}
	resp.Body.Close()
	return indyHost, true
}

func decorateChecksums(downloads []string) []string {
	downSet := make(map[string]bool)
	for _, artifact := range downloads {
		downSet[artifact] = true
		if strings.HasSuffix(artifact, ".md5") || strings.HasSuffix(artifact, ".sha1") {
			continue
		}
		downSet[artifact+".md5"] = true
		downSet[artifact+".sha1"] = true
		// downSet[artifact+".sha256"] = true
	}
	finalDownloads := []string{}
	for artifact := range downSet {
		finalDownloads = append(finalDownloads, artifact)
	}
	return finalDownloads
}

func replaceTargets(artifacts []string, oldIndyHost, targetIndyHost, buildName string) []string {
	results := []string{}
	for _, a := range artifacts {
		final := replaceTarget(a, oldIndyHost, targetIndyHost, buildName)
		results = append(results, final)
	}
	return results
}

func replaceTarget(artifact, oldIndyHost, targetIndyHost, buildName string) string {
	final := replaceHost(artifact, oldIndyHost, targetIndyHost)
	final = replaceBuildName(final, buildName)
	return final
}

func replaceHost(artifact, oldIndyHost, targetIndyHost string) string {
	// First, replace the embedded indy host to the target one
	repl := oldIndyHost
	if IsEmptyString(repl) {
		repl = artifact[strings.Index(artifact, "//")+2:]
		repl = repl[:strings.Index(repl, "/")]
	}
	return strings.ReplaceAll(artifact, repl, targetIndyHost)
}

func replaceBuildName(artifact, buildName string) string {
	// Second, if use a new build name we should replace the old one with it.
	final := artifact
	if !IsEmptyString(buildName) {
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
	}
	return final
}

// generate a random 5 digit  number for a build repo like "build-test-xxxxx"
func generateRandomBuildName() string {
	buildPrefix := "build-test-"
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 99999
	return fmt.Sprintf(buildPrefix+"%v", rand.Intn(max-min)+min)
}
