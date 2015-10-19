package main

// CLI for repo search across a number of github users
//
//
import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	//	"text/template"

	"github.com/client9/gosupplychain"
)

func main() {
	flag.Parse()

	repo := flag.Args()[0]
	parts := strings.SplitN(repo, "/", 2)

	log.Printf("%+v", parts)
	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Set GITHUB_OAUTH_TOKEN env")
	}

	gh := gosupplychain.NewGitHub(token)
	lic, err := gh.GuessLicenseFromRepo(parts[0], parts[1], "master")
	if err != nil {
		log.Fatalf("%s", err)
	}
	if lic.Type != "" {
		fmt.Printf("%s, %s\n", lic.Type, lic.File)
	}
}
