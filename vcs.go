package gosupplychain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Verbose enables verbose operation logging.
var Verbose bool

// ShowCmd controls whether VCS commands are printed.
var ShowCmd bool

// A TagCmd describes a command to list available tags
// that can be passed to Cmd.TagSyncCmd.
type TagCmd struct {
	Cmd     string // command to list tags
	Pattern string // regexp to extract tags from list
}

// Cmd is a bad abstraction around a VCS command
type Cmd struct {
	Name string
	Cmd  string // name of binary to invoke command

	CreateCmd      string   // command to download a fresh copy of a repository
	TagLookupCmd   []TagCmd // commands to lookup tags before running tagSyncCmd
	TagSyncCmd     string   // command to sync to specific tag
	TagSyncDefault string   // command to sync to default tag
	LogCmd         string   // command to list repository changelogs in an XML format
}

// vcsGit describes how to use Git.
var vcsGit = &Cmd{
	Name: "Git",
	Cmd:  "git",

	CreateCmd: "clone --depth {depth} {repo} {dir}",

	TagLookupCmd: []TagCmd{
		{"show-ref tags/{tag} origin/{tag}", `((?:tags|origin)/\S+)$`},
	},
	TagSyncCmd:     "checkout {tag}",
	TagSyncDefault: "checkout master",
	LogCmd:         "log -n {limit} --pretty={template} {rev}",
}

// vcsList lists the known version control systems
var vcsList = []*Cmd{
	vcsGit,
}

// ByCmd returns the version control system for the given
// command name (hg, git, svn, bzr).
func ByCmd(cmd string) *Cmd {
	for _, vcs := range vcsList {
		if vcs.Cmd == cmd {
			return vcs
		}
	}
	return nil
}

func expand(m map[string]string, s string) string {
	for k, v := range m {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}
	return s
}

// Create creates a new copy of repo in dir.
// The parent of dir must exist; dir must not.
func (v *Cmd) Create(dir, repo string, depth int) error {
	return v.run(".", v.CreateCmd, "dir", dir, "repo", repo, "depth", strconv.Itoa(depth))
}

// TagSync syncs the repo in dir to the named tag,
// which either is a tag returned by tags or is v.TagDefault.
// dir must be a valid VCS repo compatible with v and the tag must exist.
func (v *Cmd) TagSync(dir, tag string) error {
	if v.TagSyncCmd == "" {
		return nil
	}
	if tag != "" {
		for _, tc := range v.TagLookupCmd {
			out, err := v.runOutput(dir, tc.Cmd, "tag", tag)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`(?m-s)` + tc.Pattern)
			m := re.FindStringSubmatch(string(out))
			if len(m) > 1 {
				tag = m[1]
				break
			}
		}
	}
	if tag == "" && v.TagSyncDefault != "" {
		return v.run(dir, v.TagSyncDefault)
	}
	return v.run(dir, v.TagSyncCmd, "tag", tag)
}

// Log logs the changes for the repo in dir.
// dir must be a valid VCS repo compatible with v.
//
// WARNING: this does not issue a "download" or "sync" command.
func (v *Cmd) Log(dir string, logTemplate string, limit int) ([]byte, error) {
	return v.runOutput(dir, v.LogCmd, "limit", strconv.Itoa(limit), "template", logTemplate)
}

// LogAtRev logs the change for repo in dir at the rev revision.
// dir must be a valid VCS repo compatible with v.
// rev must be a valid revision for the repo in dir.
//
// WARNING: this does not issue a "download" or "sync" command.
//  unlike the tools/vcs
func (v *Cmd) LogAtRev(dir, rev, logTemplate string) ([]byte, error) {
	return v.runOutput(dir, v.LogCmd, "limit", strconv.Itoa(1), "template", logTemplate, "rev", rev)
}

// run runs the command line cmd in the given directory.
// keyval is a list of key, value pairs.  run expands
// instances of {key} in cmd into value, but only after
// splitting cmd into individual arguments.
// If an error occurs, run prints the command line and the
// command's combined stdout+stderr to standard error.
// Otherwise run discards the command's output.
func (v *Cmd) run(dir string, cmd string, keyval ...string) error {
	_, err := v.run1(dir, cmd, keyval, false)
	return err
}

// runVerboseOnly is like run but only generates error output to standard error in verbose mode.
func (v *Cmd) runVerboseOnly(dir string, cmd string, keyval ...string) error {
	_, err := v.run1(dir, cmd, keyval, false)
	return err
}

// runOutput is like run but returns the output of the command.
func (v *Cmd) runOutput(dir string, cmd string, keyval ...string) ([]byte, error) {
	return v.run1(dir, cmd, keyval, true)
}

// run1 is the generalized implementation of run and runOutput.
func (v *Cmd) run1(dir string, cmdline string, keyval []string, verbose bool) ([]byte, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}
	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	_, err := exec.LookPath(v.Cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"go: missing %s command. See http://golang.org/s/gogetcmd\n",
			v.Name)
		return nil, err
	}

	cmd := exec.Command(v.Cmd, args...)
	cmd.Dir = dir
	cmd.Env = envForDir(cmd.Dir)
	if ShowCmd {
		fmt.Printf("cd %s\n", dir)
		fmt.Printf("%s %s\n", v.Cmd, strings.Join(args, " "))
	}
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err = cmd.Run()
	out := buf.Bytes()
	if err != nil {
		if verbose || Verbose {
			fmt.Fprintf(os.Stderr, "# cd %s; %s %s\n", dir, v.Cmd, strings.Join(args, " "))
			os.Stderr.Write(out)
		}
		return nil, err
	}
	return out, nil
}

// envForDir returns a copy of the environment
// suitable for running in the given directory.
// The environment is the current process's environment
// but with an updated $PWD, so that an os.Getwd in the
// child will be faster.
func envForDir(dir string) []string {
	env := os.Environ()
	// Internally we only use rooted paths, so dir is rooted.
	// Even if dir is not rooted, no harm done.
	return mergeEnvLists([]string{"PWD=" + dir}, env)
}

// mergeEnvLists merges the two environment lists such that
// variables with the same name in "in" replace those in "out".
func mergeEnvLists(in, out []string) []string {
NextVar:
	for _, inkv := range in {
		k := strings.SplitAfterN(inkv, "=", 2)[0]
		for i, outkv := range out {
			if strings.HasPrefix(outkv, k) {
				out[i] = inkv
				continue NextVar
			}
		}
		out = append(out, inkv)
	}
	return out
}

// Commit contains meta data about a single commit
type Commit struct {
	Commit  string
	Author  string
	Date    string
	Message string
}

// GitLogAtRev is a special function to parse GitHub commits
// TODO clearly the CMD would be better as a interface.
func GitLogAtRev(cmd *Cmd, rootdir, rev string) ([]Commit, error) {
	line, err := cmd.LogAtRev(rootdir, rev, `{"Commit":"%H","Author":"%an <%ae>","Date":"%ad","Message":"%f"},`)
	if err != nil {
		return nil, fmt.Errorf("%s Unable to get log: %s", rootdir, err)
	}
	if len(line) < 2 {
		return nil, fmt.Errorf("%s Unable to find %q", rootdir, rev)
	}
	line = line[:len(line)-2]
	jsonbuf := make([]byte, 0, len(line)+2)
	jsonbuf = append(jsonbuf, '[')
	jsonbuf = append(jsonbuf, line...)
	jsonbuf = append(jsonbuf, ']')
	commits := []Commit{}
	err = json.Unmarshal(jsonbuf, &commits)
	if err != nil {
		return nil, fmt.Errorf("%s unable to decode: %s", rootdir, err)
	}
	return commits, nil
}
