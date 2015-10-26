package gosupplychain

import (
	"bufio"
	"encoding/json"
	"io"
	//	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
)

// A PackageError describes an error loading information about a package.
type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error
	Err         string   // the error itself
}

// Package is copy of the Package struct as listed in https://golang.org/src/cmd/go/list.go
//  oddly not exported in golang
type Package struct {
	Dir           string // directory containing package sources
	ImportPath    string // import path of package in dir
	ImportComment string // path in import comment on package statement
	Name          string // package name
	Doc           string // package documentation string
	Target        string // install path
	Shlib         string // the shared library that contains this package (only set when -linkshared)
	Goroot        bool   // is this package in the Go root?
	Standard      bool   // is this package part of the standard Go library?
	Stale         bool   // would 'go install' do anything for this package?
	Root          string // Go root or Go path dir containing this package

	// Source files
	GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles       []string // .go sources files that import "C"
	IgnoredGoFiles []string // .go sources ignored due to build constraints
	CFiles         []string // .c source files
	CXXFiles       []string // .cc, .cxx and .cpp source files
	MFiles         []string // .m source files
	HFiles         []string // .h, .hh, .hpp and .hxx source files
	SFiles         []string // .s source files
	SwigFiles      []string // .swig files
	SwigCXXFiles   []string // .swigcxx files
	SysoFiles      []string // .syso object files to add to archive

	// Cgo directives
	CgoCFLAGS    []string // cgo: flags for C compiler
	CgoCPPFLAGS  []string // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string // cgo: flags for C++ compiler
	CgoLDFLAGS   []string // cgo: flags for linker
	CgoPkgConfig []string // cgo: pkg-config names

	// Dependency information
	Imports []string // import paths used by this package
	Deps    []string // all (recursively) imported dependencies

	// Error information
	Incomplete bool            // this package or a dependency has an error
	Error      *PackageError   // error loading package
	DepsErrors []*PackageError // errors loading dependencies

	TestGoFiles  []string // _test.go files in package
	TestImports  []string // imports from TestGoFiles
	XTestGoFiles []string // _test.go files outside package
	XTestImports []string // imports from XTestGoFiles
}

// GoListPackages is a wrapper around `go list -e -json package...`
//  golang doesnt expose this in a API
//  inpsired by github.com/tools/godep which also doesnt expose this
//  as a library
func GoListPackages(name ...string) ([]Package, error) {
	if len(name) == 0 {
		return nil, nil
	}
	args := []string{"list", "-e", "-json"}
	args = append(args, name...)
	cmd := exec.Command("go", args...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	out := make([]Package, 0, 100)
	d := json.NewDecoder(r)
	for {
		info := Package{}
		err = d.Decode(&info)
		if err == io.EOF {
			break
		}
		if err != nil {
			// should never happen
			return nil, err
		}
		out = append(out, info)
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoListStd calls `go list std` to return a list of standard packages
//  This functionality is not exported programmatically.
func GoListStd() ([]string, error) {
	cmd := exec.Command("go", "list", "std")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	// there are about 148 in go1.5
	std := make([]string, 0, 200)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := string(scanner.Text())
		std = append(std, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return std, nil
}

// GoListDeps list all depedencies for the given
// list of package paths names returned in sorted order, or error
//
func GoListDeps(name ...string) ([]string, error) {
	if len(name) == 0 {
		return nil, nil
	}
	args := []string{"list", "-f", `{{  join .Deps "\n"}}`}
	args = append(args, name...)
	//	log.Printf("CMD: %v", args)
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
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
		log.Fatalf("GoListDeps Wait failed: %s", err)
	}
	paths := make([]string, 0, len(uniq))
	for k := range uniq {
		paths = append(paths, k)
	}
	sort.Strings(paths)
	return paths, nil
}