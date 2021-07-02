package promotetest

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/commonjava/indy-tests/common"
)

func Run(targetIndy, foloTrackId, promoteTargetStore string) {
	indyHost, validated := common.ValidateTargetIndy(targetIndy)
	if !validated {
		os.Exit(1)
	}

	foloTrackContent := getFoloRecord("http://"+indyHost, foloTrackId)

	b, e := json.MarshalIndent(foloTrackContent, "", "\t")
	if e == nil {
		fmt.Print(string(b))
	}
}
