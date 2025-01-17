package buildtest

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	common "github.com/commonjava/indy-tests/pkg/common"
)

const (
	TMP_DOWNLOAD_DIR = "/tmp/download"
	TMP_UPLOAD_DIR   = "/tmp/upload"
)

func Run(originalIndy, foloId, replacement, targetIndy, buildType string, processNum int) {
	origIndy := originalIndy
	if !strings.HasPrefix(origIndy, "http://") {
		origIndy = "http://" + origIndy
	}
	foloTrackContent := common.GetFoloRecord(origIndy, foloId)
	newBuildName := common.GenerateRandomBuildName()
	DoRun(originalIndy, replacement, targetIndy, buildType, newBuildName, foloTrackContent, nil, processNum, false, false)
}

// Create the repo structure and do the download/upload
func DoRun(originalIndy, replacement, targetIndy, buildType, newBuildName string, foloTrackContent common.TrackedContent,
	additionalRepos []string,
	processNum int, clearCache, dryRun bool) bool {

	common.ValidateTargetIndyOrExit(originalIndy)
	targetIndyHost, _ := common.ValidateTargetIndyOrExit(targetIndy)

	// Prepare the indy repos for the whole testing
	buildMeta := decideMeta(buildType)
	if !prepareIndyRepos("http://"+targetIndyHost, newBuildName, *buildMeta, additionalRepos, dryRun) {
		os.Exit(1)
	}

	downloadDir, uploadDir := prepareDownUploadDirectories(foloTrackContent.TrackingKey.Id, clearCache)

	downloads := prepareDownloadEntriesByFolo(targetIndy, newBuildName, foloTrackContent, additionalRepos)
	downloadFunc := func(md5str, originalArtiURL, targetArtiURL string) bool {
		fileLoc := path.Join(downloadDir, path.Base(targetArtiURL))
		if dryRun {
			fmt.Printf("Dry run download, url: %s\n", targetArtiURL)
			return true
		}
		success := common.DownloadFile(targetArtiURL, fileLoc)
		if success {
			common.Md5Check(fileLoc, md5str)
		}
		return success
	}
	broken := false
	if len(downloads) > 0 {
		fmt.Println("Start handling downloads artifacts.")
		fmt.Printf("==========================================\n\n")
		if processNum > 1 {
			broken = !concurrentRun(processNum, downloads, downloadFunc)
		} else {
			for _, down := range downloads {
				broken = !downloadFunc(down[0], down[1], down[2])
				if broken {
					break
				}
			}
		}
		fmt.Println("==========================================")
		if broken {
			fmt.Printf("Build test failed due to some downloading errors. Please see above logs to see the details.\n\n")
			os.Exit(1)
		}
		fmt.Printf("Downloads artifacts handling finished.\n\n")
	}

	uploadFunc := func(md5str, originalArtiURL, targetArtiURL string) bool {
		if dryRun {
			fmt.Printf("Dry run upload, originalArtiURL: %s, targetArtiURL: %s\n", originalArtiURL, targetArtiURL)
			return true
		}

		cacheFile := path.Join(uploadDir, path.Base(originalArtiURL))
		var downloaded bool
		if common.FileOrDirExists(cacheFile) {
			fmt.Printf("File already downloaded, reuse cacheFile: %s\n", cacheFile)
			downloaded = true
		} else {
			downloaded = common.DownloadUploadFileForCache(originalArtiURL, cacheFile)
		}
		if downloaded {
			common.Md5Check(cacheFile, md5str)
			return common.UploadFile(targetArtiURL, cacheFile)
		}
		return false
	}

	uploads := prepareUploadEntriesByFolo(originalIndy, targetIndy, newBuildName, foloTrackContent)

	if len(uploads) > 0 {
		fmt.Println("Start handling uploads artifacts.")
		fmt.Printf("==========================================\n\n")
		if processNum > 1 {
			broken = !concurrentRun(processNum, uploads, uploadFunc)
		} else {
			for _, up := range uploads {
				broken = !uploadFunc(up[0], up[1], up[2])
				if broken {
					break
				}
			}
		}
		fmt.Println("==========================================")
		if broken {
			fmt.Printf("Build test failed due to some uploadig errors. Please see above logs to see the details.\n\n")
			os.Exit(1)
		}

		fmt.Printf("Uploads artifacts handling finished.\n\n")
	}
	if !broken && !dryRun {
		targIndy := targetIndy
		if !strings.HasPrefix(targIndy, "http://") {
			targIndy = "http://" + targIndy
		}
		if common.SealFoloRecord(targIndy, newBuildName) {
			fmt.Printf("Folo record sealing succeeded for %s", newBuildName)
		} else {
			fmt.Printf("Warning: folo record sealing failed for %s", newBuildName)
		}
	}

	return true
}

// For downloads entries, we will get the paths and inject them to the final url of target indy
// as they should be directly download from target indy.
func prepareDownloadEntriesByFolo(targetIndyURL, newBuildId string, foloRecord common.TrackedContent, additionalRepos []string) map[string][]string {
	targetIndy := normIndyURL(targetIndyURL)
	result := make(map[string][]string)
	for _, down := range foloRecord.Downloads {
		var p string
		repoPath := strings.ReplaceAll(down.StoreKey, ":", "/")
		if down.AccessChannel == "GENERIC_PROXY" {
			p = path.Join("api/content", repoPath, down.Path)
		} else if common.Contains(additionalRepos, down.StoreKey) {
			p = path.Join("api/folo/track", newBuildId, repoPath, down.Path)
		} else {
			p = path.Join("api/folo/track", newBuildId, "maven/group", newBuildId, down.Path)
		}
		downUrl := fmt.Sprintf("%s%s", targetIndy, p)
		result[down.Path] = []string{down.Md5, "", downUrl}
	}
	return result
}

// For uploads entries, firstly they should be downloaded from original indy server. We use original indy server to
// make the download url, and use the target indy server to make the upload url
func prepareUploadEntriesByFolo(originalIndyURL, targetIndyURL, newBuildId string, foloRecord common.TrackedContent) map[string][]string {
	originalIndy := normIndyURL(originalIndyURL)
	targetIndy := normIndyURL(targetIndyURL)
	result := make(map[string][]string)
	for _, up := range foloRecord.Uploads {
		orgiUpUrl, targUpUrl := createUploadUrls(originalIndy, targetIndy, newBuildId, up)
		result[up.Path] = []string{up.Md5, orgiUpUrl, targUpUrl}
	}
	return result
}

func createUploadUrls(originalIndy, targetIndy, newBuildId string, up common.TrackedContentEntry) (string, string) {
	storePath := common.StoreKeyToPath(up.StoreKey) // original store, e.g, maven/hosted/build-1234
	uploadPath := path.Join("api/content", storePath, up.Path)
	orgiUpUrl := fmt.Sprintf("%s%s", originalIndy, uploadPath)                                 // original url to retrieve artifact
	alteredUploadPath := common.AlterUploadPath(up.Path, newBuildId[len(common.BUILD_TEST_):]) // replace version number
	toks := strings.Split(storePath, "/")                                                      // get package/type, e.g., maven/hosted
	targetStorePath := path.Join(toks[0], toks[1], newBuildId, alteredUploadPath)              // e.g, maven/hosted/build-913413/org/...
	targUpUrl := fmt.Sprintf("%sapi/folo/track/%s/%s", targetIndy, newBuildId, targetStorePath)
	return orgiUpUrl, targUpUrl
}

func normIndyURL(indyURL string) string {
	indy := indyURL
	if !strings.HasPrefix(indy, "http://") {
		indy = "http://" + indy
	}
	if !strings.HasSuffix(indy, "/") {
		indy = indy + "/"
	}
	return indy
}

func prepareDownUploadDirectories(buildId string, clearCache bool) (string, string) {
	// use "/tmp/download", which will be dropped after each run
	downloadDir := TMP_DOWNLOAD_DIR
	if !common.FileOrDirExists(downloadDir) {
		os.MkdirAll(downloadDir, os.FileMode(0755))
	}
	if !common.FileOrDirExists(downloadDir) {
		fmt.Printf("Error: cannot create directory %s for file downloading.\n", downloadDir)
		os.Exit(1)
	}

	// use ENVAR_TEST_MOUNT_PATH + "bulidId/upload" if this envar is defined
	uploadDir := TMP_UPLOAD_DIR
	envarTestMountPath := os.Getenv(common.ENVAR_TEST_MOUNT_PATH)
	if envarTestMountPath != "" {
		uploadDir = path.Join(envarTestMountPath, buildId, "upload")
		if clearCache {
			os.RemoveAll(uploadDir)
		}
	}

	if !common.FileOrDirExists(uploadDir) {
		os.MkdirAll(uploadDir, os.FileMode(0755))
	}

	if !common.FileOrDirExists(uploadDir) {
		fmt.Printf("Error: cannot create directory %s for caching uploading files.\n", uploadDir)
		os.Exit(1)
	}
	fmt.Printf("Prepared download dir: %s, upload dir: %s\n", downloadDir, uploadDir)
	return downloadDir, uploadDir
}

func concurrentRun(numWorkers int, artifacts map[string][]string, job func(md5, originalURL, targetURL string) bool) bool {
	fmt.Printf("Start to run job in concurrent mode with thread number %v\n", numWorkers)
	ch := make(chan []string, numWorkers*5) // This buffered number of chan can be anything as long as it's larger than numWorkers
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results = []bool{}

	// This starts numWorkers number of goroutines that wait for something to do
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				mu.Lock()
				results = append(results, job(a[0], a[1], a[2]))
				mu.Unlock()
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for _, artifact := range artifacts {
		ch <- artifact // add artifact to the queue
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish

	finalResult := true
	for _, result := range results {
		if finalResult = result; !finalResult {
			break
		}
	}

	return finalResult
}
