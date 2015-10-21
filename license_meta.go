package gosupplychain

// LicenseMeta is struct containing various meta data about a license
// including names and links to other websites.
type LicenseMeta struct {
	FullName      string // Full name in English
	LinkOriginal  string // Link to original license source
	LinkOSI       string // Link to The Open Source Initiative, http://opensource.org
	LinkOSIAlt    string // Alternate link for The Open Source Initiative (normally old)
	LinkCAL       string // Link to "Choose a License"
	LinkTLDR      string // Link to "TLDR;Legal"
	LinkWikipedia string // Link to Wikipedia
}

// Meta is a mapping from license tokens to meta data.
var Meta = map[string]LicenseMeta{
	"Apache-2.0": {
		FullName:      "Apache License 2.0",
		LinkOriginal:  "http://www.apache.org/licenses/license-2.0",
		LinkOSI:       "http://opensource.org/licenses/Apache-2.0",
		LinkCAL:       "http://choosealicense.com/licenses/apache-2.0/",
		LinkTLDR:      "https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/Apache_License",
	},
	"NewBSD": {
		FullName:      "BSD 3-Clause License",
		LinkOSI:       "http://opensource.org/licenses/BSD-3-Clause",
		LinkCAL:       "http://choosealicense.com/licenses/bsd-3-clause/",
		LinkTLDR:      "https://tldrlegal.com/license/bsd-3-clause-license-(revised)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/BSD_licenses",
	},
	"FreeBSD": {
		FullName:      "BSD 2-Clause License",
		LinkOSI:       "http://opensource.org/licenses/BSD-2-Clause",
		LinkCAL:       "http://choosealicense.com/licenses/bsd-2-clause/",
		LinkTLDR:      "https://tldrlegal.com/license/bsd-2-clause-license-(freebsd)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/BSD_licenses",
	},
	"GPL-2.0": {
		FullName:      "GNU General Public License v2",
		LinkOriginal:  "http://www.gnu.org/licenses/old-licenses/gpl-2.0.html",
		LinkOSI:       "http://opensource.org/licenses/GPL-2.0",
		LinkCAL:       "http://choosealicense.com/licenses/gpl-2.0/",
		LinkTLDR:      "https://tldrlegal.com/license/gnu-general-public-license-v2",
		LinkWikipedia: "https://en.wikipedia.org/wiki/GNU_General_Public_License",
	},
	"GPL-3.0": {
		FullName:      "GNU General Public License v3",
		LinkOriginal:  "http://www.gnu.org/licenses/gpl.html",
		LinkOSI:       "http://opensource.org/licenses/GPL-3.0",
		LinkCAL:       "http://choosealicense.com/licenses/gpl-3.0/",
		LinkTLDR:      "https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/GNU_General_Public_License",
	},
	"LGPL-2.1": {
		FullName:      "GNU Lesser General Public License v2.1",
		LinkOriginal:  "http://www.gnu.org/licenses/old-licenses/lgpl-2.1.html",
		LinkOSI:       "http://opensource.org/licenses/LGPL-2.1",
		LinkCAL:       "http://choosealicense.com/licenses/lgpl-2.1/",
		LinkTLDR:      "https://tldrlegal.com/license/gnu-lesser-general-public-license-v2.1-(lgpl-2.1)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/GNU_Lesser_General_Public_License",
	},
	"LGPL-3.0": {
		FullName:      "GNU Lesser General Public License v3.0",
		LinkOriginal:  "http://www.gnu.org/licenses/lgpl.html",
		LinkOSI:       "http://opensource.org/licenses/LGPL-3.0",
		LinkCAL:       "http://choosealicense.com/licenses/lgpl-3.0/",
		LinkTLDR:      "https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/GNU_Lesser_General_Public_License",
	},
	"MIT": {
		FullName:      "MIT License",
		LinkOSI:       "http://opensource.org/licenses/MIT",
		LinkCAL:       "http://choosealicense.com/licenses/mit/",
		LinkTLDR:      "https://tldrlegal.com/license/mit-license",
		LinkWikipedia: "https://en.wikipedia.org/wiki/MIT_License",
	},
	"MPL-2.0": {
		FullName:      "Mozilla Public License 2.0",
		LinkOriginal:  "https://www.mozilla.org/en-US/MPL/2.0/",
		LinkOSI:       "http://opensource.org/licenses/MPL-2.0",
		LinkCAL:       "http://choosealicense.com/licenses/mpl-2.0/",
		LinkTLDR:      "https://tldrlegal.com/license/mozilla-public-license-2.0-(mpl-2)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/Mozilla_Public_License",
	},
	"AGPL-3.0": {
		FullName:      "GNU Affero General Public License",
		LinkOriginal:  "http://www.gnu.org/licenses/agpl.html",
		LinkOSI:       "http://opensource.org/licenses/AGPL-3.0",
		LinkCAL:       "http://choosealicense.com/licenses/agpl-3.0/",
		LinkWikipedia: "https://en.wikipedia.org/wiki/Affero_General_Public_License",
	},
	"WTFPL-2.0": {
		FullName:      "Do What The Fuck You Want To Public License",
		LinkOriginal:  "http://www.wtfpl.net/txt/copying/",
		LinkTLDR:      "https://tldrlegal.com/license/do-wtf-you-want-to-public-license-v2-(wtfpl-2.0)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/WTFPL",
	},
	"CDDL-1.0": {
		FullName:      "Common Development and Distribution License",
		LinkTLDR:      "https://tldrlegal.com/license/common-development-and-distribution-license-(cddl-1.0)-explained",
		LinkWikipedia: "https://en.wikipedia.org/wiki/Common_Development_and_Distribution_License",
	},
	"EPL-1.0": {
		FullName:      "Eclipse Public License 1.0",
		LinkOriginal:  "https://www.eclipse.org/legal/epl-v10.html",
		LinkTLDR:      "https://tldrlegal.com/license/eclipse-public-license-1.0-(epl-1.0)",
		LinkWikipedia: "https://en.wikipedia.org/wiki/Eclipse_Public_License",
	},
	"Unlicense": {
		FullName:     "Unlicense",
		LinkOriginal: "http://unlicense.org",
	},
}
