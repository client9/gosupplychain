package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	//	"github.com/client9/go-license"
	"github.com/client9/gosupplychain"
	"github.com/client9/gosupplychain/golist"
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

	t, err := template.New("internal").Funcs(golist.TemplateFuncMap()).Parse(bomTemplate)
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
