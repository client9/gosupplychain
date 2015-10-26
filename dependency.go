package gosupplychain

//import ("github.com/golang/gddo/gosrc"
//	"net/http"
//)
import (
	"log"

	"net/mail"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/client9/gosupplychain/golist"

	"golang.org/x/tools/go/vcs"

	"github.com/ryanuber/go-license"
)

// Notes:
//  go-source meta tag:  https://github.com/golang/gddo/wiki/Source-Code-Links

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
	log.Printf("Directory is %s", dir)
	cmd := exec.Command("git", "log", "-1", "--format=Commit: %H%nDate: %aD%nSubject: %s%n%n%b%n")
	cmd.Dir = dir
	msg, err := cmd.Output()
	if err != nil {
		log.Printf("git log error: %s", err)
		return Commit{}, err
	}
	//log.Printf("GOT %s", string(msg))
	r := strings.NewReader(string(msg))
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Printf("git log parse error: %s", err)
		return Commit{}, err
	}
	header := m.Header

	//behind, _ := GitCommitsBehind(dir, header.Get("Commit"))

	return Commit{
		Date:    header.Get("Date"),
		Hash:    header.Get("Commit"),
		Subject: header.Get("Subject"),
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
func LoadDependencies(pkgs []string, ignores []string) ([]ExternalDependency, error) {

	stdlib, err := golist.Std()
	if err != nil {
		return nil, err
	}

	pkgs, err = golist.Deps(pkgs...)
	if err != nil {
		return nil, err
	}

	// faster to remove stdlib
	pkgs = removeIfEquals(pkgs, stdlib)
	pkgs = removeIfSubstring(pkgs, ignores)

	deps, err := golist.Packages(pkgs...)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool, len(deps))

	out := make([]ExternalDependency, 0, len(deps))
	for _, v := range deps {
		src := filepath.Join(v.Root, "src")
		path := filepath.Join(src, filepath.FromSlash(v.ImportPath))
		cmd, root, err := vcs.FromDir(path, src)
		if err != nil {
			log.Printf("  failed %s", err)
		} else {
			log.Printf("  got root %s  cmd=%s", root, cmd.Cmd)
		}

		rr, err := vcs.RepoRootForImportPath(v.ImportPath, false)
		if err != nil {
			log.Printf("   unable to get repo root for %s: %s", v.ImportPath, err)
		} else {
			log.Printf("   repo = %s root=%s", rr.Repo, rr.Root)
		}

		visited[v.ImportPath] = true

		e := ExternalDependency{
			Package: v.ImportPath,
		}
		l, err := license.NewFromDir(path)
		if err == nil {
			e.File = filepath.Base(l.File)
			e.License = l.Type
		} else if !visited[rr.Root] {
			visited[rr.Root] = true
			path = filepath.Join(src, filepath.FromSlash(rr.Root))
			l, err = license.NewFromDir(path)
			if err == nil {
				e.Package = rr.Root
				e.File = filepath.Base(l.File)
				e.License = l.Type
			}
		} else {
			continue
		}
		c, err := GetLastCommit(path)
		if err == nil {
			e.Commit = c.Hash
			e.Date = c.Date
		}

		if strings.HasPrefix(e.Package, "github.com") && e.Commit != "" && e.File != "" {
			e.LicenseLink = "https://" + e.Package + "/blob/" + e.Commit + "/" + e.File
		} else if strings.HasPrefix(e.Package, "golang.org/x/") && e.Commit != "" && e.File != "" {
			e.LicenseLink = "https://github.com/golang/" + e.Package[13:] + "/blob/" + e.Commit + "/" + e.File
		} else if strings.HasPrefix(e.Package, "gopkg.in") && e.Commit != "" && e.File != "" {
			e.LicenseLink = GoPkgInToGitHub(e.Package) + "/" + e.File
		}
		out = append(out, e)
	}
	return out, err
}

// generic []string function, remove elements of A that are in B
func removeIfEquals(alist []string, blist []string) []string {
	out := make([]string, 0, len(alist))
	for _, a := range alist {
		add := true
		for _, b := range blist {
			if a == b {
				add = false
			}
		}
		if add {
			out = append(out, a)
		}
	}
	return out
}

// removes elements of A that substring match any of B
func removeIfSubstring(alist []string, blist []string) []string {
	out := make([]string, 0, len(alist))
	for _, a := range alist {
		add := true
		for _, b := range blist {
			if strings.Index(a, b) != -1 {
				add = false
			}
		}
		if add {
			out = append(out, a)
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
func addParents(pkgs []string) []string {
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
