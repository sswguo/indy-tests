package buildtest

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	common "github.com/commonjava/indy-tests/common"
)

const TMP_DOWNLOAD_DIR = "/tmp/download"
const TMP_UPLOAD_DIR = "/tmp/upload"

func Run(indyURL, foloId, replacement, targetIndy, buildType string, processNum int) {
	indyHost, validated := common.ValidateTargetIndy(targetIndy)
	if !validated {
		os.Exit(1)
	}

	newBuildName := generateRandomBuildName()

	// Prepare the indy repos for the whole testing
	buildMeta := decideMeta(buildType)
	if !prepareIndyRepos("http://"+indyHost, newBuildName, *buildMeta) {
		os.Exit(1)
	}

	// result := prepareEntriesByLog(logUrl)

	result := prepareEntriesByFolo(indyURL, foloId)

	prepareCacheDirectories()

	fmt.Println("Start handling downloads artifacts.")
	fmt.Printf("==========================================\n\n")
	downloads := replaceTargets(result["downloads"], "", indyHost, newBuildName)
	result["downloads"] = nil // save memory
	downloadFunc := func(artifactUrl string) {
		fileLoc := path.Join(TMP_DOWNLOAD_DIR, path.Base(artifactUrl))
		common.DownloadFile(artifactUrl, fileLoc)
	}
	if downloads != nil {
		if processNum > 1 {
			concurrentRun(processNum, downloads, downloadFunc)
		} else {
			for _, url := range downloads {
				downloadFunc(url)
			}
		}
	}
	fmt.Println("==========================================")
	fmt.Printf("Downloads artifacts handling finished.\n\n")

	fmt.Println("Start handling uploads artifacts.")
	fmt.Printf("==========================================\n\n")
	// uploads := replaceTargets(result["uploads"], "", indyHost, newBuildName)
	// result["uploads"] = nil // save memory
	uploadFunc := func(artifactUrl string) {
		cacheFile := path.Join(TMP_UPLOAD_DIR, path.Base(artifactUrl))
		downloadArtifact := replaceHost(artifactUrl, "", indyHost)
		downloaded := common.DownloadUploadFileForCache(downloadArtifact, cacheFile)
		if downloaded {
			replacedUrl := replaceBuildName(downloadArtifact, newBuildName)
			common.UploadFile(replacedUrl, cacheFile)
		}
	}
	if result["uploads"] != nil {
		if processNum > 1 {
			concurrentRun(processNum, result["uploads"], uploadFunc)
		} else {
			for _, url := range result["uploads"] {
				uploadFunc(url)
			}
		}
	}
	fmt.Println("==========================================")
	fmt.Printf("Uploads artifacts handling finished.\n\n")
}

func prepareEntriesByFolo(indyURL, foloId string) map[string][]string {
	indy := indyURL
	if !strings.HasPrefix(indy, "http://") {
		indy = "http://" + indy
	}
	foloTrackContent := common.GetFoloRecord(indy, foloId)
	indyFinalURL := indy
	if !strings.HasSuffix(indyFinalURL, "/") {
		indyFinalURL = indyFinalURL + "/"
	}
	result := make(map[string][]string)
	downloads := []string{}
	for _, down := range foloTrackContent.Downloads {
		downUrl := fmt.Sprintf("%sapi/folo/track/%s/maven/group/%s%s", indyFinalURL, foloId, foloId, down.Path)
		downloads = append(downloads, downUrl)
	}
	result["downloads"] = downloads
	uploads := []string{}
	for _, up := range foloTrackContent.Uploads {
		upUrl := fmt.Sprintf("%sapi/folo/track/%s/maven/group/%s%s", indyFinalURL, foloId, foloId, up.Path)
		uploads = append(uploads, upUrl)
	}
	result["uploads"] = uploads
	return result
}

func prepareCacheDirectories() {
	if !common.FileOrDirExists(TMP_DOWNLOAD_DIR) {
		os.Mkdir(TMP_DOWNLOAD_DIR, os.FileMode(0755))
	}
	if !common.FileOrDirExists(TMP_DOWNLOAD_DIR) {
		fmt.Printf("Error: cannot create directory %s for file downloading.\n", TMP_DOWNLOAD_DIR)
		os.Exit(1)
	}
	if !common.FileOrDirExists(TMP_UPLOAD_DIR) {
		os.Mkdir(TMP_UPLOAD_DIR, os.FileMode(0755))
	}
	if !common.FileOrDirExists(TMP_UPLOAD_DIR) {
		fmt.Printf("Error: cannot create directory %s for caching uploading files.\n", TMP_UPLOAD_DIR)
		os.Exit(1)
	}
}

// Deprecated as folo is the preferred way now.
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
	if common.IsEmptyString(repl) {
		repl = artifact[strings.Index(artifact, "//")+2:]
		repl = repl[:strings.Index(repl, "/")]
	}
	return strings.ReplaceAll(artifact, repl, targetIndyHost)
}

func replaceBuildName(artifact, buildName string) string {
	// Second, if use a new build name we should replace the old one with it.
	final := artifact
	if !common.IsEmptyString(buildName) {
		buildPat := regexp.MustCompile(`https{0,1}:\/\/.+\/(build-\S+?)\/.*`)
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

func concurrentRun(numWorkers int, artifacts []string, job func(artifact string)) {
	var ch = make(chan string, numWorkers*5) // This buffered number of chan can be anything as long as it's larger than numWorkers
	var wg sync.WaitGroup

	// This starts xthreads number of goroutines that wait for something to do
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				job(a) // do the thing
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for _, artifact := range artifacts {
		ch <- artifact // add i to the queue
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish
}
