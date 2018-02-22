package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/types"
	"io"
	"log"
	"strings"

	"github.com/podhmo/astknife/action"
	"github.com/podhmo/astknife/bypos"
	"github.com/podhmo/astknife/lookup"
	"github.com/podhmo/strangejson/buildcontext"
	"github.com/podhmo/strangejson/output"
	"github.com/podhmo/strangejson/output/codegen/fileinfo"
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
	cache := &cache{
		unmarshalStructCache: map[*types.Struct]types.Object{},
	}

	for _, pkgpath := range pkgpaths {
		info := cmd.Program.Package(pkgpath)
		fset := cmd.Program.Fset
		sorted := bypos.SortFiles(info.Files)
		lookup := lookup.New(sorted.Files...)

		gens, err := GenerateUnmarshalJSON(info.Pkg, sorted, cache.unmarshalStructCache)
		if err != nil {
			return err
		}

		r := fileinfo.NewRepository(fset, info.Pkg, info.Files)
		aggregated := map[*ast.File][]Gen{}
		for _, gen := range gens {
			log.Printf("for %s.%s\n", gen.Object.Type().String(), gen.Name)
			srcfile := gen.File
			srcfilename, _ := r.NameOf(srcfile)

			dstfilename := strings.Replace(strings.Replace(srcfilename, ".go", "_gen.go", 1), "_gen_gen.go", "_gen.go", 1)
			dstfile, ok := r.FileOf(dstfilename)

			if !ok {
				dstfile = r.CreateFakeFile(dstfilename)
				r.AddFile(dstfilename, dstfile)
			}
			aggregated[dstfile] = append(aggregated[dstfile], gen)
		}

		for f, gens := range aggregated {
			filename, _ := r.NameOf(f)
			log.Printf("write %s\n", filename)

			var newCodeBuf bytes.Buffer
			io.WriteString(&newCodeBuf, fmt.Sprintf("package %s\n", info.Pkg.Name()))
			for _, gen := range gens {
				if err := gen.Generate(&newCodeBuf, f); err != nil {
					return err
				}
			}

			switch r.IsNew(f) {
			case true:
				output, err := imports.Process(filename, newCodeBuf.Bytes(), nil)
				if err != nil {
					return err
				}
				buildcontext.WriteFile(cmd.Build, filename, output, 0644)
			default:
				newf, err := parser.ParseFile(fset, filename, newCodeBuf.Bytes(), parser.ParseComments)
				if err != nil {
					return err
				}

				for _, gen := range gens {
					replacement := lookup.With(newf).Lookup(gen.Name) // xxx
					if _, err := action.AppendOrReplace(lookup, f, replacement); err != nil {
						if !action.IsNoEffect(err) {
							return err
						}
					}
				}

				var b bytes.Buffer
				if err := printer.Fprint(&b, fset, f); err != nil {
					return err
				}

				output, err := imports.Process(filename, b.Bytes(), nil)
				if err != nil {
					return err
				}
				buildcontext.WriteFile(cmd.Build, filename, output, 0644)
			}
		}
	}
	return nil
}

// Gen :
type Gen struct {
	Name     string
	Object   types.Object
	File     *ast.File
	Generate func(w io.Writer, f *ast.File) error // xxx tentative
}
