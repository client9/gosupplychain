package gosupplychain

import (
	"log"
	"strings"

	"golang.org/x/tools/go/vcs"
)

type CommitMini struct {
	SHA string
	Msg string
}

type ImportStatus struct {
	Root    string       // root import
	Status  string       // ahead, or behind
	Commits []CommitMini // specific
}

func cleanText(msg string) string {
	msg = strings.Replace(msg, "\t", " ", -1)
	msg = strings.Replace(msg, "\r", " ", -1)
	msg = strings.Replace(msg, "\n", " ", -1)
	msg = strings.Replace(msg, "  ", " ", -1)
	msg = strings.TrimSpace(msg)
	if len(msg) > 80 {
		msg = msg[:80] + "..."
	}
	return msg
}

// Behind takes a github token and a godep file
//  and returns a list of dependencies and if they are out of date
func Behind(githubToken string, godepFile string) []ImportStatus {
	gh := NewGitHub(githubToken)
	gd, err := LoadGodepsFile(godepFile)
	if err != nil {
		log.Fatalf("Error loading godeps file %q: %s", godepFile, err)
	}

	roots := make(map[string]bool, len(gd.VendorDeps()))

	imports := make([]ImportStatus, 0, len(gd.VendorDeps()))

	for _, dep := range gd.VendorDeps() {
		rr, err := vcs.RepoRootForImportPath(dep.Name, false)
		if err != nil {
			log.Printf("Unable to process %s: %s", dep.Name, err)
			continue
		}
		if roots[rr.Root] {
			continue
		}
		roots[rr.Root] = true
		parts := strings.Split(dep.Name, "/")
		if len(parts) < 2 {
			log.Printf("Skipping %s", dep.Name)
			continue
		}
		if parts[0] == "golang.org" && parts[1] == "x" {
			parts[0] = "github.com"
			parts[1] = "golang"
		}

		if parts[0] != "github.com" {
			log.Printf("Skipping %s", dep.Name)
			continue
		}

		compare, _, err := gh.Client.Repositories.CompareCommits(parts[1], parts[2], dep.Version, "HEAD")
		if err != nil {
			log.Printf("got error reading repo %s: %s", dep.Name, err)
			continue
		}
		status := ImportStatus{
			Root:    rr.Root,
			Status:  *compare.Status,
			Commits: make([]CommitMini, 0, len(compare.Commits)),
		}
		for _, c := range compare.Commits {
			msg := ""
			if c.Commit.Message != nil {
				msg = *c.Commit.Message
			}
			status.Commits = append(status.Commits, CommitMini{
				SHA: *c.SHA,
				Msg: msg,
			})
		}
		imports = append(imports, status)
	}
	return imports
}
