package promotetest

import (
	"fmt"
	"os"

	"github.com/commonjava/indy-tests/common"
)

func Run(targetIndy, foloTrackId, promoteTargetStore string) {
	indyHost, validated := common.ValidateTargetIndy(targetIndy)
	if !validated {
		os.Exit(1)
	}

	indyURL := "http://" + indyHost

	foloTrackContent := getFoloRecord(indyURL, foloTrackId)

	if foloTrackContent.Uploads == nil && len(foloTrackContent.Uploads) == 0 {
		fmt.Printf("There are not any uploads records in folo build %s, promotion will be ignored!\n", foloTrackId)
		return
	}

	paths := []string{}

	sourcePromote := foloTrackContent.Uploads[0].StoreKey
	for _, entry := range foloTrackContent.Uploads {
		paths = append(paths, entry.Path)
	}

	promote(indyURL, sourcePromote, promoteTargetStore, paths)
}
