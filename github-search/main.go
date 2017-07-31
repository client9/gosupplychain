package main

// CLI for repo search across a number of github users
//
//
import (
	"flag"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/client9/gosupplychain"
	"github.com/ryanuber/go-license"
)

// SearchMarkdownTemplate is a sample markdown output
var SearchMarkdownTemplate = `
{{ range $index, $user := $.Users }}
### {{ $user.Name }}

| Package | Description | License | Updated |
|---------|-------------|---------|---------|
{{ range $user.Repos }}| [{{ .Name }}](https://github.com/{{ .Name }}) | {{ .Description }} | {{ with (index $.Licenses .Name) }} [{{ .Type }}]({{ .File }}){{ end }} | {{ .Updated.Format "2006-01-02" }} |
{{ end }}
{{ end }}
`

func main() {
	searchQuery := flag.String("query", "language:go", "Search query to be executed per user")
	addLicense := flag.Bool("add-license", true, "Attempt to determine software license (slower)")

	// TODO add flag for template
	// TODO add flag for output file

	flag.Parse()

	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Set GITHUB_OAUTH_TOKEN env")
	}

	names := flag.Args()
	sort.Strings(names)

	gh := gosupplychain.NewGitHub(token)
	users, err := gh.SearchByUsers(token, *searchQuery, names)
	if err != nil {
		log.Fatalf("Github failed: %s", err)
	}

	licmap := make(map[string]license.License)
	if *addLicense {
		for _, user := range users {
			for _, repo := range user.Repos {
				log.Printf("Checking license for %s", repo.Name)
				parts := strings.SplitN(repo.Name, "/", 2)
				lic, err := gh.GuessLicenseFromRepo(parts[0], parts[1], "master")
				if err != nil {
					log.Printf("Unable to check license for %s: %s", repo.Name, err)
				}
				log.Printf("... Got %s", lic.Type)
				licmap[repo.Name] = lic
			}
		}
	}
	t, err := template.New("test").Parse(SearchMarkdownTemplate)
	if err != nil {
		log.Fatalf("Template init failed: %s", err)
	}

	err = t.Execute(os.Stdout, map[string]interface{}{
		"Users":    users,
		"Licenses": licmap,
	})
	if err != nil {
		log.Fatalf("Template exec failed: %s", err)
	}
}
