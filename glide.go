package gosupplychain

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// A GlideDependency is a specific revision of a package.
type GlideDependency struct {
	Name    string
	Version string
}

// Glide describes a glide.lock file
type Glide struct {
	Imports []GlideDependency
}

func (g Glide) VendorDeps() []VendorDependency {
	var deps []VendorDependency
	for _, dep := range g.Imports {
		deps = append(deps, VendorDependency{
			Name:    dep.Name,
			Version: dep.Version,
		})
	}

	return deps
}

// LoadGlideFile loads a glide file
func LoadGlideFile(path string) (PackageManager, error) {
	var g Glide
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return g, err
	}
	err = yaml.Unmarshal(b, &g)
	return g, err
}
