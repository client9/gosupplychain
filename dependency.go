package gosupplychain

//import ("github.com/golang/gddo/gosrc"
//	"net/http"
//)
import (
	"log"

	"net/mail"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/client9/gosupplychain/golist"

	"golang.org/x/tools/go/vcs"

	"github.com/ryanuber/go-license"
)

// Notes:
//  go-source meta tag:  https://github.com/golang/gddo/wiki/Source-Code-Links

// Project contains an amalgamation of package, commit, repo, and license information
type Project struct {
	VcsName     string
	VcsCmd      string
	Repo        string
	LicenseLink string
}

// Commit contains meta data about a single commit
type Commit struct {
	Rev     string
	Date    string
	Subject string
	Body    string
	Behind  int
}

// Dependency contains meta data on a external dependency
type Dependency struct {
	golist.Package
	Commit  Commit
	License license.License
	Project Project
}

// LinkToFile returns a URL that links to particular revision of a
// file or empty
//
func LinkToFile(pkg, file, rev string) string {
	if file == "" {
		return ""
	}

	switch {
	case strings.HasPrefix(pkg, "github.com"):
		if rev == "" {
			rev = "master"
		}
		return "https://" + pkg + "/blob/" + rev + "/" + file
	case strings.HasPrefix(pkg, "golang.org/x/"):
		if rev == "" {
			rev = "master"
		}
		return "https://github.com/golang/" + pkg[13:] + "/blob/" + rev + "/" + file
	case strings.HasPrefix(pkg, "gopkg.in"):
		return GoPkgInToGitHub(pkg) + "/" + file
	default:
		return ""
	}
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
		Rev:     header.Get("Commit"),
		Subject: header.Get("Subject"),
		Behind:  -1,
	}, nil
}

// GetLicense returns licensing info
func GetLicense(path string) license.License {

	l, err := license.NewFromDir(path)
	if err != nil {
		return license.License{}
	}
	l.File = filepath.Base(l.File)
	return *l
}

// LoadDependencies is not done
func LoadDependencies(pkgs []string, ignores []string) ([]Dependency, error) {

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

	out := make([]Dependency, 0, len(deps))
	for _, v := range deps {
		src := filepath.Join(v.Root, "src")
		path := filepath.Join(src, filepath.FromSlash(v.ImportPath))
		cmd, _, err := vcs.FromDir(path, src)
		rr, err := vcs.RepoRootForImportPath(v.ImportPath, false)
		visited[v.ImportPath] = true

		e := Dependency{
			Package: v,
		}
		e.Project.Repo = rr.Repo
		e.Project.VcsName = cmd.Name
		e.Project.VcsCmd = cmd.Cmd
		e.License = GetLicense(path)
		if e.License.Type == "" && !visited[rr.Root] {
			visited[rr.Root] = true

			// BUG: really need to call go-list again to get info
			path = filepath.Join(src, filepath.FromSlash(rr.Root))
			e.License = GetLicense(path)
		}
		commit, err := GetLastCommit(path)
		if err == nil {
			e.Commit = commit
		}

		e.Project.LicenseLink = LinkToFile(e.ImportPath, e.License.File, e.Commit.Rev)

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
