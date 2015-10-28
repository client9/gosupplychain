package gosupplychain

import (
	"encoding/json"
	"os"
)

// A GoDepDependency is a specific revision of a package.
type GoDepDependency struct {
	ImportPath string
	Comment    string `json:",omitempty"` // Description of commit, if present.
	Rev        string // VCS-specific commit ID.
}

// Godeps describes what a package needs to be rebuilt reproducibly.
// It's the same information stored in file Godeps.
type Godeps struct {
	ImportPath string
	GoVersion  string
	Packages   []string `json:",omitempty"` // Arguments to save, if any.
	Deps       []GoDepDependency
}

// LoadGodepsFile loads a godeps file
func LoadGodepsFile(path string) (Godeps, error) {
	var g Godeps
	f, err := os.Open(path)
	if err != nil {
		return g, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&g)
	return g, err
}
