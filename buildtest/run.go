package buildtest

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func Run(logUrl, replacement, targetIndy string) {
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

	BASE_DIR := "/tmp/download"
	if !FileOrDirExists(BASE_DIR) {
		os.Mkdir(BASE_DIR, os.FileMode(0755))
	}
	if !FileOrDirExists(BASE_DIR) {
		fmt.Printf("Error: cannot create directory %s for file downloading.\n", BASE_DIR)
		os.Exit(1)
	}
	if err == nil {
		downloads := replaceTarget(decorateChecksums(result["downloads"]), "", targetIndy)
		result["downloads"] = nil
		if downloads != nil {
			for _, url := range downloads {
				fileLoc := path.Join(BASE_DIR, path.Base(url))
				DownloadFile(url, fileLoc)
			}
		}
		uploads := replaceTarget(result["uploads"], "", targetIndy)
		result["uploads"] = nil
		if uploads != nil {
			fmt.Print("Start showing uploads: ==================\n\n")
			for _, u := range uploads {
				fmt.Println(u)
			}
			fmt.Print("\nFinish showing uploads: ==================\n\n")
		}
	}
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

func replaceTarget(artifacts []string, replacement, target string) []string {
	results := []string{}
	for _, a := range artifacts {
		repl := replacement
		if IsEmptyString(repl) {
			repl = a[strings.Index(a, "//")+2:]
			repl = repl[:strings.Index(repl, "/")]
		}
		final := strings.ReplaceAll(a, repl, target)
		// final = strings.ReplaceAll(final, "folo/track/build-xxxx/maven/group/build-xxxx", "content/maven/hosted/shared-imports")
		fmt.Println(final)
		results = append(results, final)
	}
	return results
}
