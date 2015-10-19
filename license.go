package gosupplychain

import (
	"strings"
)

// LicenseFilePrefix is a list of filename prefixes that indicate it
//  might contain a software license
var LicenseFilePrefix = []string{
	"license",
	"copying",
	"unlicense",
	"copyright",
}

// IsPossibleLicenseFile returns true if the filename might be contain a software license
func IsPossibleLicenseFile(filename string) bool {
	lowerfile := strings.ToLower(filename)
	for _, prefix := range LicenseFilePrefix {
		if strings.HasPrefix(lowerfile, prefix) {
			return true
		}
	}
	return false
}
