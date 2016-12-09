package gosupplychain

// VendorDependency is used as a common depedency structure
// between package management tools
type VendorDependency struct {
	Name    string
	Version string
}

// PackageManager is a package management
// tool (Godeps, glide, etc.)
type PackageManager interface {
	VendorDeps() []VendorDependency
}
