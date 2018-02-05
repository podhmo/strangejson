package swaggergen

import (
	"github.com/k0kubun/pp"
	"github.com/podhmo/strangejson/output"
	"golang.org/x/tools/go/loader"
)

type command struct {
	Program *loader.Program
}

// New :
func New(program *loader.Program) output.Command {
	return &command{Program: program}
}

// run :
func (cmd *command) Run(pkgpaths []string) error {
	for _, pkgpath := range pkgpaths {
		info := cmd.Program.Package(pkgpath)
		findDescription := true
		schemas, err := ParsePackageInfo(info, findDescription)
		if err != nil {
			return err
		}
		pp.Println(schemas)
	}
	return nil
}
