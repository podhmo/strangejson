package formatchecktask

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/astknife/bypos"
	"github.com/podhmo/strangejson/output"
	"github.com/podhmo/strangejson/output/codegen/accessor"
	"github.com/podhmo/strangejson/output/codegen/typesutil"
)

// Task :
type Task struct {
	formatcheckable *FormatCheckable
}

// New :
func New(formatcheckable *FormatCheckable) output.Task {
	return &Task{
		formatcheckable: formatcheckable,
	}
}

// Prepare :
func (t *Task) Prepare(pkg *types.Package, files bypos.Sorted) error {
	a := &accessor.Accessor{Pkg: pkg}
	return a.IterateStructs(func(sa *accessor.StructAccessor) error {
		t.formatcheckable.RegisterFake(sa.Object.Type())
		return nil
	})
}

// Do :
func (t *Task) Do(pkg *types.Package, files bypos.Sorted) ([]output.Gen, error) {
	a := &accessor.Accessor{Pkg: pkg}
	var results []output.Gen

	err := a.IterateStructs(func(sa *accessor.StructAccessor) error {
		if !sa.Exported() {
			return nil
		}

		results = append(results, output.Gen{
			Name:   fmt.Sprintf("%s.FormatCheck", sa.Name()),
			Object: sa.Object,
			File:   bypos.FindFile(files, sa.Object.Pos()),
			Generate: func(w io.Writer, f *ast.File) error {
				qf := typesutil.ImportedNameTo(pkg, f)
				return generateFormatCheck(w, f, a, sa, qf, t.formatcheckable)
			},
		})
		return nil
	})

	return results, err
}
