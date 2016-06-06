package main

// CLI for determining how ahead/behind a repo is from master
import (
	//	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/client9/gosupplychain"
)

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

	imports := gosupplychain.Behind(token, args[0])

	for _, imp := range imports {
		if imp.Status == "identical" {
			continue
		}
		fmt.Printf("%s: %s\n", imp.Root, imp.Status)
		for i, c := range imp.Commits {
			fmt.Printf("  %d: %s %s\n", i, c.SHA[0:7], cleanText(c.Msg))
		}
	}

	/*
		raw, err := json.MarshalIndent(imports, "", "  ")
		if err != nil {
			log.Fatalf("unable to marshal output: %s", err)
		}
		fmt.Printf("%s\n", raw)
	*/
}
