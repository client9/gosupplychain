package main

import (
	"flag"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/client9/gosupplychain"
)

// SearchMarkdownTemplate is a sample markdown output
var SearchMarkdownTemplate = `
{{ range $index, $user := $ }}
### {{ $user.Name }}

| Package | Description | Updated |
|---------|-------------|---------|
{{ range $user.Repos }}| [{{ .Name }}](https://github.com/{{ .Name }}) | {{ .Description }} | {{ .Updated.Format "2006-01-02" }} |
{{ end }}
{{ end }}
`

func main() {
	searchQuery := flag.String("query", "language:go", "Search query to be executed per user")

	// TODO add flag for template
	// TODO add flag for output file

	flag.Parse()

	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Set GITHUB_OAUTH_TOKEN env")
	}

	names := flag.Args()
	sort.Strings(names)

	users, err := gosupplychain.GitHubSearchByUsers(token, *searchQuery, names)
	if err != nil {
		log.Fatalf("Github failed: %s", err)
	}

	t, err := template.New("test").Parse(SearchMarkdownTemplate)
	if err != nil {
		log.Fatalf("Template init failed: %s", err)
	}

	err = t.Execute(os.Stdout, users)
	if err != nil {
		log.Fatalf("Template exec failed: %s", err)
	}
}
