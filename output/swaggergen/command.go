package swaggergen

import (
	"github.com/k0kubun/pp"
	"github.com/podhmo/strangejson/buildcontext"
	"github.com/podhmo/strangejson/output"
	"golang.org/x/tools/go/loader"
)

type command struct {
	Build   *buildcontext.Context
	Program *loader.Program
}

// New :
func New(build *buildcontext.Context, program *loader.Program) output.Command {
	return &command{Build: build, Program: program}
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
