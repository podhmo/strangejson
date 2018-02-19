package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"io"
	"log"
	"strings"

	"go/printer"

	"github.com/podhmo/astknife/bypos"
	"github.com/podhmo/strangejson/buildcontext"
	"github.com/podhmo/strangejson/output"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/imports"
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
		fset := cmd.Program.Fset
		sorted := bypos.SortFiles(info.Files)

		gens, err := GenerateUnmarshalJSON(info.Pkg)
		if err != nil {
			return err
		}

		filenameMap := make(map[*ast.File]string, len(info.Files))
		fileMap := map[string]*ast.File{}
		for _, f := range info.Files {
			filename := fset.Position(f.Pos()).Filename
			filenameMap[f] = filename
			fileMap[filename] = f
		}

		aggregated := map[*ast.File][]Gen{}
		for _, gen := range gens {
			log.Printf("for %s.%s\n", gen.Object.Type().String(), gen.Name)
			srcfile := bypos.FindFile(sorted, gen.Object.Pos())
			srcfilename := filenameMap[srcfile]

			dstfilename := strings.Replace(strings.Replace(srcfilename, ".go", "_gen.go", 1), "_gen_gen.go", "_gen.go", 1)
			dstfile, ok := fileMap[dstfilename]

			if !ok {
				code := fmt.Sprintf("package %s", info.Pkg.Name())
				dstfile, err = parser.ParseFile(fset, dstfilename, code, parser.ParseComments)
				if err != nil {
					return err
				}
				fileMap[dstfilename] = dstfile
				filenameMap[dstfile] = dstfilename
			}
			aggregated[dstfile] = append(aggregated[dstfile], gen)
		}

		for f, gens := range aggregated {
			filename := filenameMap[f]
			log.Printf("write %s\n", filename)

			// todo: string -> ast. sync ast.
			var b bytes.Buffer
			if err := printer.Fprint(&b, fset, f); err != nil {
				return err
			}
			for _, gen := range gens {
				if err := gen.Generate(&b, f); err != nil {
					return err
				}
			}

			output, err := imports.Process(filename, b.Bytes(), nil)
			if err != nil {
				return err
			}

			buildcontext.WriteFile(cmd.Build, filename, output, 0644)
		}
	}
	return nil
}

// Gen :
type Gen struct {
	Name     string
	Object   types.Object
	Generate func(w io.Writer, f *ast.File) error // xxx tentative
}
