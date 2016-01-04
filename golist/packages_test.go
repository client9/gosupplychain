package golist

import (
	//	"log"
	"testing"
)

// TestGoListStd tests GoListStd.  This is really over kill
//  for the current implementation, but previous one was much weirder
func TestGoListStd(t *testing.T) {
	pkgs, err := Std()
	if err != nil {
		t.Fatalf("Unable to get standard packages: %s", err)
	}
	if len(pkgs) == 0 {
		t.Fatalf("No packages found!")
	}
	pmap := make(map[string]bool)
	for _, pkgs := range pkgs {
		pmap[pkgs] = true
	}
	cases := []struct {
		path string
		has  bool
	}{
		{"archive", false},
		{"archive/tar", true},
		{"bytes", true},
		{"go/internal", false},
		{"compress/bzip2/testdata", false},
		{"text/template/parse", true},
		{"unsafe", true},
	}

	for pos, tt := range cases {
		_, found := pmap[tt.path]
		if found != tt.has {
			if found {
				t.Errorf("case %d: path %q is not a standard package", pos, tt.path)
			} else {
				t.Errorf("case %d: path %q is a standard package", pos, tt.path)
			}
		}
	}
}

func TestNewContext(t *testing.T) {
	c, err := NewContext()
	if err != nil {
		t.Fatalf("Unable to get context: %s", err)
	}
	if c.GOARCH != "" {
		t.Errorf("GOARCH not set correctly: %q", c.GOARCH)
	}
}
