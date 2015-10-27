package gosupplychain

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// MetaGoImport represents the values in a go-import meta tag.
type MetaGoImport struct {
	ProjectRoot string
	Vcs         string
	Repo        string
}

// MetaGoSource represents the values in a go-source meta tag.
type MetaGoSource struct {
	ProjectRoot  string
	ProjectURL   string
	DirTemplate  string
	FileTemplate string
}

// DirURL returns a URL pointing to the VCS directory
func (mgs MetaGoSource) DirURL(dir string) string {
	return replaceDir(mgs.DirTemplate, dir)
}

// FileURL returns a URL points to the VCS File
func (mgs MetaGoSource) FileURL(dir, file string) string {
	tpl := replaceDir(mgs.FileTemplate, dir)
	parts := strings.SplitN(tpl, "#", 2)
	return strings.Replace(parts[0], "{file}", file, -1)
}

func replaceDir(s string, dir string) string {
	slashDir := ""
	dir = strings.Trim(dir, "/")
	if dir != "" {
		slashDir = "/" + dir
	}
	s = strings.Replace(s, "{dir}", dir, -1)
	s = strings.Replace(s, "{/dir}", slashDir, -1)
	return s
}

// charsetReader returns a reader for the given charset. Currently
// it only supports UTF-8 and ASCII. Otherwise, it returns a meaningful
// error which is printed by go get, so the user can find why the package
// wasn't downloaded if the encoding is not supported. Note that, in
// order to reduce potential errors, ASCII is treated as UTF-8 (i.e. characters
// greater than 0x7f are not rejected).
func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "ascii":
		return input, nil
	default:
		return nil, fmt.Errorf("can't decode XML document using charset %q", charset)
	}
}

// parseMetaGoImports returns meta imports from the HTML in r.
// Parsing ends at the end of the <head> section or the beginning of the <body>.
func parseMetaGo(r io.Reader) (mgi *MetaGoImport, mgs *MetaGoSource, err error) {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	d.Strict = false
	var t xml.Token
	for {
		t, err = d.RawToken()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if e, ok := t.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, "body") {
			return
		}
		if e, ok := t.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, "head") {
			return
		}
		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}
		switch attrValue(e.Attr, "name") {
		case "go-import":
			if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 3 {
				mgi = &MetaGoImport{
					ProjectRoot: f[0],
					Vcs:         f[1],
					Repo:        f[2],
				}
			}
		case "go-source":
			if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 4 {
				mgs = &MetaGoSource{
					ProjectRoot:  f[0],
					ProjectURL:   f[1],
					DirTemplate:  f[2],
					FileTemplate: f[3],
				}
			}
		}
	}
}

// attrValue returns the attribute value for the case-insensitive key
// `name', or the empty string if nothing is found.
func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}
