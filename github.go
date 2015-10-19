package gosupplychain

import (
	"fmt"
	"log"
	"time"

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

// GitHubSearchByUsers performs a search on multiple users
func GitHubSearchByUsers(oauthToken string, searchQuery string, users []string) ([]User, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: oauthToken,
		})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
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
	for _, co := range users {
		q := fmt.Sprintf("user:%s %s", co, searchQuery)
		log.Printf("Running query %q", q)
		repos, _, err := client.Search.Repositories(q, opts)
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
		time.Sleep(time.Second * 4)
	}
	return out, nil
}
