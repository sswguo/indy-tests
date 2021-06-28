package promotetest

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	common "github.com/commonjava/indy-tests/common"
)

// IndyPromoteVars ...
type IndyPromoteVars struct {
	Source string
	Target string
	Paths  []string
}

// IndyPromoteJSONTemplate ...
func IndyPromoteJSONTemplate(indyPromoteVars *IndyPromoteVars) string {
	request := `{
  "async": false,
  "source": "{{.Source}}",
  "target": "{{.Target}}",
  {{if gt (len .Paths) 0}}
  "paths": [{{range $index,$path := .Paths}}"{{$path}}"{{if isNotLast $index $.Paths}},{{end}}{{end}}],
  {{end}}
  "purgeSource": false,
  "dryRun": false,
  "fireEvents": true,
  "failWhenExists": true
}`

	t := template.Must(template.New("settings").Funcs(isNotLast).Parse(request))
	var buf bytes.Buffer
	err := t.Execute(&buf, indyPromoteVars)
	if err != nil {
		log.Fatal("executing template:", err)
		os.Exit(1)
	}

	return buf.String()
}

var isNotLast = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"isNotLast": func(index int, array []string) bool {
		return index < len(array)-1
	},
}

func promote(indyURL, source, target string, paths []string) {
	promoteVars := IndyPromoteVars{
		Source: source,
		Target: target,
		Paths:  paths,
	}
	promote := IndyPromoteJSONTemplate(&promoteVars)

	URL := fmt.Sprintf("%s/api/promotion/paths/promote", indyURL)

	fmt.Printf("Start promote request:\n %s\n\n", promote)
	respText, result := common.PostRequest(URL, strings.NewReader(promote))

	if result {
		fmt.Printf("Promote Done. Result is:\n %s\n\n", respText)
	} else {
		fmt.Printf("Promote Error. Result is:\n %s\n\n", respText)
	}
}
