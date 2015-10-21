package gosupplychain

// LicenseMeta is struct containing various meta data about a license
// including names and links to other websites.
type LicenseMeta struct {
	FullName     string // Full name in English
	LinkOriginal string // Link to original license source
	LinkOSI      string // Link to The Open Source Initiative, http://opensource.org
	LinkOSIAlt   string // Alternate link for The Open Source Initiative
	LinkCAL      string // Link to "Choose a License"
	LinkTLDR     string // Link to "TLDR;Legal"
}

// Meta is a mapping from license tokens to meta data.
var Meta = map[string]LicenseMeta{
	"Apache-2.0": {
		FullName:     "Apache License 2.0",
		LinkOriginal: "www.apache.org/licenses/license-2.0",
		LinkOSI:      "http://opensource.org/licenses/Apache-2.0",
		LinkCAL:      "http://choosealicense.com/licenses/apache-2.0/",
		LinkTLDR:     "https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)",
	},
	"NewBSD": {
		FullName: "BSD 3-Clause License",
		LinkOSI:  "http://opensource.org/licenses/BSD-3-Clause",
		LinkCAL:  "http://choosealicense.com/licenses/bsd-3-clause/",
		LinkTLDR: "https://tldrlegal.com/license/bsd-3-clause-license-(revised)",
	},
	"FreeBSD": {
		FullName: "BSD 2-Clause License",
		LinkOSI:  "http://opensource.org/licenses/BSD-2-Clause",
		LinkCAL:  "http://choosealicense.com/licenses/bsd-2-clause/",
		LinkTLDR: "https://tldrlegal.com/license/bsd-2-clause-license-(freebsd)",
	},
	"GPL-2.0": {
		FullName: "GNU General Public License v2",
		LinkOSI:  "http://opensource.org/licenses/GPL-2.0",
		LinkCAL:  "http://choosealicense.com/licenses/gpl-2.0/",
		LinkTLDR: "https://tldrlegal.com/license/gnu-general-public-license-v2",
	},
	"GPL-3.0": {
		FullName: "GNU General Public License v3",
		LinkOSI:  "http://opensource.org/licenses/GPL-3.0",
		LinkCAL:  "http://choosealicense.com/licenses/gpl-3.0/",
		LinkTLDR: "https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)",
	},
	"LGPL-2.1": {
		FullName: "GNU Lesser General Public License v2.1",
		LinkOSI:  "http://opensource.org/licenses/LGPL-2.1",
		LinkCAL:  "http://choosealicense.com/licenses/lgpl-2.1/",
		LinkTLDR: "https://tldrlegal.com/license/gnu-lesser-general-public-license-v2.1-(lgpl-2.1)",
	},
	"LGPL-3.0": {
		FullName: "GNU Lesser General Public License v3.0",
		LinkOSI:  "http://opensource.org/licenses/LGPL-3.0",
		LinkCAL:  "http://choosealicense.com/licenses/lgpl-3.0/",
		LinkTLDR: "https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)",
	},
	"MIT": {
		FullName: "MIT License",
		LinkOSI:  "http://opensource.org/licenses/MIT",
		LinkCAL:  "http://choosealicense.com/licenses/mit/",
		LinkTLDR: "https://tldrlegal.com/license/mit-license",
	},
	"MPL-2.0": {
		FullName: "Mozilla Public License 2.0",
		LinkOSI:  "http://opensource.org/licenses/MPL-2.0",
		LinkCAL:  "http://choosealicense.com/licenses/mpl-2.0/",
		LinkTLDR: "https://tldrlegal.com/license/mozilla-public-license-2.0-(mpl-2)",
	},
	"AGPL-3.0": {
		FullName: "XXX",
		LinkOSI:  "http://opensource.org/licenses/AGPL-3.0",
		LinkCAL:  "http://choosealicense.com/licenses/agpl-3.0/",
	},
	"WTFPL-2.0": {
		FullName: "Do What The Fuck You Want To Public License",
		LinkTLDR: "https://tldrlegal.com/license/do-wtf-you-want-to-public-license-v2-(wtfpl-2.0)",
	},
	"CDDL-1.0": {
		FullName: "Common Development and Distribution License",
		LinkTLDR: "https://tldrlegal.com/license/common-development-and-distribution-license-(cddl-1.0)-explained",
	},
	"EPL-1.0": {
		FullName: "Eclipse Public License 1.0",
		LinkTLDR: "https://tldrlegal.com/license/eclipse-public-license-1.0-(epl-1.0)",
	},
	"Unlicense": {
		FullName: "XXX",
	},
}
