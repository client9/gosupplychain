package gosupplychain

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/client9/go-license"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Repo describes a repo basic
// NOTE: likely to be replaced with a larger structure
type Repo struct {
	Name        string
	Description string
	Updated     time.Time
}

// User is the top level GitHub user (maybe be a company or user)
// NOTE: like to be replaced with a larger structure
type User struct {
	Name  string
	Repos []Repo
}

// GitHubFile is contains everything needed to represent a file at a point in time
//  Likely to be generalized later
type GitHubFile struct {
	Owner string
	Repo  string
	Path  string
	Tree  string
	SHA   string
}

// RawURL returns a URL to the raw content, without formatting
func (file GitHubFile) RawURL() string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", file.Owner, file.Repo, file.Tree, file.Path)
}

// WebURL returns a human-friend URL to github
func (file GitHubFile) WebURL() string {
	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", file.Owner, file.Repo, file.Tree, file.Path)
}

// GitHub is a VCS
type GitHub struct {
	Client *github.Client
}

// NewGitHub creates a github client using oauth token
func NewGitHub(oauthToken string) GitHub {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: oauthToken,
		})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return GitHub{
		Client: github.NewClient(tc),
	}
}

// GetFileContentsURL generates a download URL
func (gh GitHub) GetFileContentsURL(owner, repo, sha, filepath string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, sha, filepath)
}

// GetFileContents down loads a file
func (gh GitHub) GetFileContents(owner, repo, tree, filepath string) (string, error) {
	url := gh.GetFileContentsURL(owner, repo, tree, filepath)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

// GetTreeFiles returns the list of files given a tree.
//
// sha must be a valid git sha value or "master"
func (gh GitHub) GetTreeFiles(owner string, repo string, sha string) ([]GitHubFile, error) {
	tree, _, err := gh.Client.Git.GetTree(owner, repo, sha, false)
	if err != nil {
		return nil, err
	}
	//log.Printf("TREE: %+v", *tree)
	out := make([]GitHubFile, 0, len(tree.Entries))
	for _, t := range tree.Entries {
		out = append(out, GitHubFile{
			Owner: owner,
			Repo:  repo,
			Tree:  sha,
			Path:  *t.Path,
			SHA:   *t.SHA,
		})

		//log.Printf("TREE: %s", t)
	}
	return out, nil
}

// SearchByUsers performs a search on multiple users
func (gh GitHub) SearchByUsers(oauthToken string, searchQuery string, users []string) ([]User, error) {
	opts := &github.SearchOptions{
		Sort:  "updated",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	out := make([]User, 0, len(users))

	// assume each query takes 1 second round trip
	//  and we get 20/minute
	//  wait 4 seconds between calls
	for pos, co := range users {
		if pos > 0 {
			time.Sleep(time.Second * 4)
		}
		q := fmt.Sprintf("user:%s %s", co, searchQuery)
		log.Printf("Running query %q", q)
		repos, _, err := gh.Client.Search.Repositories(q, opts)
		if err != nil {
			return nil, err
		}
		if repos == nil || *repos.Total == 0 || repos.Repositories == nil {
			continue
		}
		user := User{
			Name: co,
		}
		for _, val := range repos.Repositories {

			r := Repo{}
			if val.FullName != nil {
				r.Name = *val.FullName
			}
			if val.Description != nil {
				r.Description = *val.Description
			}
			if val.UpdatedAt != nil {
				// UpdateAt is a odd github.Time that embeds a time.Time
				tmp := *val.UpdatedAt
				r.Updated = tmp.Time
			}
			user.Repos = append(user.Repos, r)
		}
		out = append(out, user)
	}
	return out, nil
}

// GuessLicenseFromRepo attempts to determine a license
func (gh GitHub) GuessLicenseFromRepo(owner string, repo string, sha string) (license.License, error) {

	files, err := gh.GetTreeFiles(owner, repo, sha)
	if err != nil {
		return license.License{}, err
	}
	out := []string{}
	for _, filename := range files {
		if IsLicenseFile(filename.Path) {
			out = append(out, filename.Path)
			body, err := gh.GetFileContents(owner, repo, sha, filename.Path)
			if err != nil {
				return license.License{}, fmt.Errorf("unable to download %s: %s", filename, err)
			}
			lic := license.License{
				Text: body,
				File: filename.WebURL(),
			}
			err = lic.GuessType()
			if err == nil {
				return lic, nil
			}
		}
	}
	// empty
	return license.License{}, nil
}
