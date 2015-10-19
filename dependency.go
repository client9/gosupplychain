package gosupplychain

import (
	"bufio"
	//	"fmt"
	"log"
	"net/mail"

	//	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ryanuber/go-license"
)

// ExternalDependency contains meta data on a external dependency
type ExternalDependency struct {
	Package     string
	File        string
	License     string
	LicenseLink string
	Commit      string
	Date        string
	Behind      int
}

// Commit contains meta data about a single commit
type Commit struct {
	Hash    string
	Date    string
	Subject string
	Body    string
	Behind  int
}

// GoPkgInToGitHub converts a "gopkg.in" to a github repo link
func GoPkgInToGitHub(name string) string {
	parts := strings.Split(name, "/")
	if parts[0] != "gopkg.in" {
		return ""
	}
	pname := ""
	version := ""
	if strings.Index(parts[1], ".") != -1 {
		versionparts := strings.SplitN(parts[1], ".", 2)
		pname = versionparts[0]
		version = versionparts[1]
		return "https://github.com/go-" + pname + "/" + pname + "/blob/" + version + "/" + parts[2]
	}
	versionparts := strings.SplitN(parts[len(parts)-1], ".", 2)
	pname = versionparts[0]
	version = versionparts[1]
	return "https://github.com/" + parts[1] + "/" + pname + "/tree/" + version
}

// GitFetch executes a "git fetch" command
func GitFetch(dir string) error {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = dir
	return cmd.Run()
}

// GitCommitsBehind counts the number of commits a directory is behind master
func GitCommitsBehind(dir string, hash string) (int, error) {
	// the following doesnt work sometimes
	//cmd := exec.Command("git", "rev-list", "..master")
	cmd := exec.Command("git", "rev-list", "--count", "origin/master..."+hash)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(strings.TrimSpace(string(out)))
}

// GetLastCommit returns meta data on the last commit
func GetLastCommit(dir string) (Commit, error) {
	/*
		err := GitFetch(dir)
		if err != nil {
			log.Fatalf("unable to fetch: %s", err)
		}
	*/
	//log.Printf("Directory is %s", dir)
	cmd := exec.Command("git", "log", "-1", "--format=Commit: %H%nDate: %aD%nSubject: %s%n%n%b%n")
	cmd.Dir = dir
	msg, err := cmd.Output()
	if err != nil {
		return Commit{}, err
	}
	//log.Printf("GOT %s", string(msg))
	r := strings.NewReader(string(msg))
	m, err := mail.ReadMessage(r)
	if err != nil {
		return Commit{}, err
	}
	header := m.Header

	//GitFetch(dir)
	behind, _ := GitCommitsBehind(dir, header.Get("Commit"))

	return Commit{
		Date:    header.Get("Date"),
		Hash:    header.Get("Commit"),
		Subject: header.Get("Subject"),
		Behind:  behind,
	}, nil
}

// GetLicense returns licensing info
func GetLicense(rootpath, name string) ExternalDependency {
	sw := ExternalDependency{
		Package: name,
		File:    "",
		License: "",
	}

	l, err := license.NewFromDir(filepath.Join(rootpath, "src", name))
	if err != nil {
		//log.Printf(err.Error())
		return sw
	}

	sw.File = filepath.Base(l.File)
	sw.License = l.Type
	return sw
}

// LoadDependencies is not done
func LoadDependencies(root string, pkgs []string, ignores []string) []ExternalDependency {

	pkgs, err := ListDependenciesFromPackages(pkgs...)
	if err != nil {
		log.Printf("FAILED: %s", err)
		return nil
	}

	pkgs = RemoveStandardPackages(pkgs)
	pkgs = RemoveIgnores(pkgs, ignores)
	pkgs = AddParents(pkgs)

	externals := make([]ExternalDependency, 0, 100)

	for _, v := range pkgs {
		e := GetLicense(root, v)
		if e.License == "" && len(strings.Split(v, "/")) > 3 {
			continue
		}

		commit, err := GetLastCommit(filepath.Join(root, "src", v))
		if err == nil {
			e.Commit = commit.Hash
			e.Date = commit.Date
			e.Behind = commit.Behind
		}

		if strings.HasPrefix(e.Package, "github.com") && e.Commit != "" && e.File != "" {
			e.LicenseLink = "https://" + e.Package + "/blob/" + e.Commit + "/" + e.File
		} else if strings.HasPrefix(e.Package, "golang.org/x/") && e.Commit != "" && e.File != "" {
			name := filepath.Base(e.Package)
			e.LicenseLink = "https://go.googlesource.com/" + name + "/+/" + e.Commit + "/" + e.File
		} else if strings.HasPrefix(e.Package, "gopkg.in") && e.Commit != "" && e.File != "" {
			e.LicenseLink = GoPkgInToGitHub(e.Package) + "/" + e.File
		}

		externals = append(externals, e)

	}
	return externals
}

// ListDependenciesFromPackages list all depedencies for the given list of paths
//  returns in sorted order, or error
func ListDependenciesFromPackages(name ...string) ([]string, error) {
	if len(name) == 0 {
		return nil, nil
	}
	args := []string{"list", "-f", `{{  join .Deps "\n"}}`}
	args = append(args, name...)
	//	log.Printf("CMD: %v", args)
	cmd := exec.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	uniq := make(map[string]bool)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		uniq[scanner.Text()] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Wait failed: %s", err)
	}
	paths := make([]string, 0, len(uniq))
	for k := range uniq {
		paths = append(paths, k)
	}
	sort.Strings(paths)
	return paths, nil
}

// RemoveStandardPackages removes standard packages using a heuristic
func RemoveStandardPackages(pkgs []string) []string {
	out := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		// has dots.. not part of standard library
		if strings.Index(pkg, ".") != -1 {
			out = append(out, pkg)
		}
	}
	return out
}

// RemoveIgnores removes packages we dont want included
func RemoveIgnores(pkgs []string, ignores []string) []string {
	out := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		add := true
		for _, pattern := range ignores {
			if strings.Index(pkg, pattern) != -1 {
				add = false
			}
		}
		if add {
			out = append(out, pkg)
		}
	}
	return out
}

// AddParents adds parent depedencies
//
// e.g. given
//   golang.org/x/crypto/hmac
//
// then the following are added
//   golang.org/x/crypto
//   golang.org/x/crypto/hmac
//
// Sometimes licenses are in the child directory, or the parent directory
//
func AddParents(pkgs []string) []string {
	uniq := make(map[string]bool, len(pkgs))
	for _, pkg := range pkgs {
		// ignore items with 2 parts, and add all subdirectories
		parts := strings.Split(pkg, "/")
		for i := 3; i <= len(parts); i++ {
			newpath := strings.Join(parts[0:i], "/")
			uniq[newpath] = true
		}
	}
	paths := make([]string, 0, len(uniq))
	for k := range uniq {
		paths = append(paths, k)
	}
	sort.Strings(paths)
	return paths
}
