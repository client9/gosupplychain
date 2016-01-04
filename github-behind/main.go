package main

// CLI for determining how ahead/behind a repo is from master
import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/client9/gosupplychain"
	//	"github.com/google/go-github/github"
)

func main() {
	flag.Parse()

	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Set GITHUB_OAUTH_TOKEN env")
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Need path to godeps file")
	}
	gd, err := gosupplychain.LoadGodepsFile(args[0])
	if err != nil {
		log.Fatalf("Error loading godeps file %q: %s", args[0], err)
	}
	gh := gosupplychain.NewGitHub(token)

	for _, dep := range gd.Deps {
		parts := strings.Split(dep.ImportPath, "/")
		if len(parts) < 2 {
			log.Printf("Skipping %s", dep.ImportPath)
			continue
		}
		if parts[0] == "golang.org" && parts[1] == "x" {
			parts[0] = "github.com"
			parts[1] = "golang"
		}

		if parts[0] != "github.com" {
			log.Printf("Skipping %s", dep.ImportPath)
			continue
		}

		compare, _, err := gh.Client.Repositories.CompareCommits(parts[1], parts[2], dep.Rev, "HEAD")
		if err != nil {
			log.Printf("got error reading repo %s: %s", dep.ImportPath, err)
			continue
		}

		fmt.Printf("%s: %s\n", dep.ImportPath, *compare.Status)
		for pos, commit := range compare.Commits {
			msg := ""
			if commit.Commit.Message != nil {
				msg = *commit.Commit.Message
				msg = strings.Replace(msg, "\t", " ", -1)
				msg = strings.Replace(msg, "\r", " ", -1)
				msg = strings.Replace(msg, "\n", " ", -1)
				msg = strings.Replace(msg, "  ", " ", -1)
				if len(msg) > 80 {
					msg = msg[:80] + "..."
				}
			}
			sha := *commit.SHA
			fmt.Printf("    %d %s %s\n", pos, sha[0:7], msg)
		}
	}
}
