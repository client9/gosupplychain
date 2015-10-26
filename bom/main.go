package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	//	"github.com/client9/go-license"
	"github.com/client9/gosupplychain"
)

var bomTemplate = `
{{ range .Depends }}
{{ .Package }}
{{ .File }}
{{ .Date }}
{{ .Commit }}
{{ .License }} {{ (index $.License .License).FullName }}
{{ .LicenseLink }}
{{ end }}
`

var bomTemplate2 = `
{{ range .Depends }}
NAME: {{ .Name }}
IMPORT: {{ .ImportPath }}
ROOT {{ .Root }}
DOC {{ .Doc }}
{{ end }}
`

func main() {
	ignoreFlag := flag.String("ignore", "", "Comma-separated stuff to skip")
	flag.Parse()
	ignores := []string{"internal"}
	if *ignoreFlag != "" {
		ignores = append(ignores, strings.Split(*ignoreFlag, ",")...)
	}
	deps, err := gosupplychain.LoadDependencies(flag.Args(), ignores)
	if err != nil {
		log.Fatalf("Unable to load dependencies: %s", err)
	}

	t, err := template.New("test").Parse(bomTemplate)
	if err != nil {
		log.Fatalf("Template init failed: %s", err)
	}

	err = t.Execute(os.Stdout, map[string]interface{}{
		"Depends": deps,
		"License": gosupplychain.Meta,
	})
	if err != nil {
		log.Fatalf("Template exec failed: %s", err)
	}
}
