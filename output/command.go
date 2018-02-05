package output

// Command :
type Command interface {
	Run(pkgpaths []string) error
}
