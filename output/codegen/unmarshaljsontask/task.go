package unmarshaljsontask

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
	arrived map[*types.Struct]types.Object
}

// New :
func New(arrived map[*types.Struct]types.Object) output.Task {
	return &Task{
		arrived: arrived,
	}
}

// Prepare :
func (t *Task) Prepare(pkg *types.Package, files bypos.Sorted) error {
	return nil
}

// Do :
func (t *Task) Do(pkg *types.Package, files bypos.Sorted) ([]output.Gen, error) {
	a := &accessor.Accessor{Pkg: pkg}
	var results []output.Gen

	err := a.IterateStructs(func(sa *accessor.StructAccessor) error {
		if !sa.Exported() {
			return nil
		}

		sameOb, ok := t.arrived[sa.Underlying]
		if !ok {
			t.arrived[sa.Underlying] = sa.Object
		}
		results = append(results, output.Gen{
			Name:   fmt.Sprintf("%s.UnmarshalJSON", sa.Name()),
			Object: sa.Object,
			File:   bypos.FindFile(files, sa.Object.Pos()),
			Generate: func(w io.Writer, f *ast.File) error {
				qf := typesutil.ImportedNameTo(pkg, f)
				return generateUnmarshalJSON(w, f, a, sa, qf, sameOb)
			},
		})
		return nil
	})

	return results, err
}
