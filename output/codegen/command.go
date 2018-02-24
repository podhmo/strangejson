package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/podhmo/astknife/bypos"
	"github.com/podhmo/astknife/lookup"
	"github.com/podhmo/strangejson/buildcontext"
	"github.com/podhmo/strangejson/output"
	"github.com/podhmo/strangejson/output/codegen/formatchecktask"
	"github.com/podhmo/strangejson/output/codegen/unmarshaljsontask"
	"golang.org/x/tools/go/ast/astutil"
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

type target struct {
	info   *loader.PackageInfo
	sorted bypos.Sorted
}

// run :
func (cmd *command) Run(pkgpaths []string) error {
	formatcheckable, err := formatchecktask.NewFormatCheckable(cmd.Program)
	if err != nil {
		return err
	}

	unmarshalStructCache := map[*types.Struct]types.Object{}

	var targets []target
	for _, pkgpath := range pkgpaths {
		info := cmd.Program.Package(pkgpath)
		targets = append(targets, target{
			info:   info,
			sorted: bypos.SortFiles(info.Files),
		})
	}

	tasks := map[*loader.PackageInfo][]output.Task{}
	for _, target := range targets {
		{
			t := unmarshaljsontask.New(unmarshalStructCache)
			if err := t.Prepare(target.info.Pkg, target.sorted); err != nil {
				return err
			}
			tasks[target.info] = append(tasks[target.info], t)
		}

		{
			t := formatchecktask.New(formatcheckable)
			if err := t.Prepare(target.info.Pkg, target.sorted); err != nil {
				return err
			}
			tasks[target.info] = append(tasks[target.info], t)
		}
	}

	for _, target := range targets {
		fset := cmd.Program.Fset

		var gens []output.Gen
		for _, t := range tasks[target.info] {
			newgens, err := t.Do(target.info.Pkg, target.sorted)
			if err != nil {
				return err
			}
			gens = append(gens, newgens...)
		}

		r := NewRepository(fset, target.info.Pkg, target.info.Files)

		aggregated := map[*ast.File][]output.Gen{}
		for _, gen := range gens {
			log.Printf("generate %s (for %s)\n", gen.Name, gen.Object.Type().String())
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

		lookup := lookup.New(target.sorted.Files...)
		for f, gens := range aggregated {
			filename, _ := r.NameOf(f)
			log.Printf("write %s\n", filename)

			var newCodeBuf bytes.Buffer
			io.WriteString(&newCodeBuf, fmt.Sprintf("package %s\n", target.info.Pkg.Name()))

			// xxx tentative
			io.WriteString(&newCodeBuf, "import (\n")
			io.WriteString(&newCodeBuf, "	\"github.com/pkg/errors\"")
			io.WriteString(&newCodeBuf, ")\n")

			sort.Slice(gens, func(i, j int) bool { return gens[i].Name < gens[j].Name })
			for _, gen := range gens {
				if err := gen.Generate(&newCodeBuf, f); err != nil {
					return err
				}
			}

			switch r.IsNew(f) {
			case true:
				if err := emitFile(cmd.Build, filename, newCodeBuf.Bytes()); err != nil {
					return err
				}
			default:
				astutil.RewriteImport(fset, f, "errors", "github.com/pkg/errors")
				modifier := &Modifier{lookup: lookup, fset: fset, f: f}
				code, err := modifier.modifyCode(filename, gens, newCodeBuf.Bytes())
				if err != nil {
					return err
				}
				if err := emitFile(cmd.Build, filename, code); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func emitFile(ctxt *buildcontext.Context, filename string, body []byte) error {
	output, err := imports.Process(filename, body, nil)
	if err != nil {
		return err
	}
	return buildcontext.WriteFile(ctxt, filename, output, 0644)
}
