package codegen

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"sort"

	"github.com/podhmo/astknife/bypos"
	"github.com/podhmo/strangejson/output/codegen/accessor"
)

// GenerateUnmarshalJSON :
func GenerateUnmarshalJSON(pkg *types.Package, sorted bypos.Sorted, arrived map[*types.Struct]types.Object) ([]Gen, error) {
	a := &accessor.Accessor{Pkg: pkg}
	var results []Gen

	err := a.IterateStructs(func(sa *accessor.StructAccessor) error {
		if !sa.Exported() {
			return nil
		}

		sameOb, ok := arrived[sa.Underlying]
		if !ok {
			arrived[sa.Underlying] = sa.Object
		}
		results = append(results, Gen{
			Name:   fmt.Sprintf("%s.UnmarshalJSON", sa.Name()),
			Object: sa.Object,
			File:   bypos.FindFile(sorted, sa.Object.Pos()),
			Generate: func(w io.Writer, f *ast.File) error {
				qf := NameTo(pkg, f)
				return generateUnmarshalJSON(w, f, a, sa, qf, sameOb)
			},
		})
		return nil
	})

	sort.Slice(results, func(i, j int) bool { return results[i].Name < results[j].Name })

	return results, err
}

func generateUnmarshalJSON(w io.Writer, f *ast.File, a *accessor.Accessor, sa *accessor.StructAccessor, qf types.Qualifier, sameOb types.Object) error {
	ob := sa.Object

	fmt.Fprintf(w, "// UnmarshalJSON : (generated from %s)\n", ob.Type().String())
	fmt.Fprintf(w, "func (x %s) UnmarshalJSON(b []byte) error {\n", types.TypeString(types.NewPointer(ob.Type()), qf))
	defer fmt.Fprintln(w, "}")

	if sameOb != nil {
		fmt.Fprintf(w, "	return (%s)(x).UnmarshalJSON(b)\n", types.TypeString(types.NewPointer(sameOb.Type()), qf))
		return nil
	}

	// internal struct, all fields are pointer
	fmt.Fprintf(w, "	type internal struct {\n")
	sa.IterateFields(func(fa *accessor.FieldAccessor) error {
		if !fa.Exported() {
			return nil
		}
		typ := types.TypeString(types.NewPointer(fa.Object.Type()), qf)
		switch fa.Tag {
		case "":
			fmt.Fprintf(w, "		%s %s\n", fa.Name(), typ)
		default:
			fmt.Fprintf(w, "		%s %s `%s`\n", fa.Name(), typ, fa.Tag)
		}
		return nil
	})
	fmt.Fprintln(w, "	}")

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "	var p internal")
	fmt.Fprintln(w, "	if err := json.Unmarshal(b, &p); err != nil{")
	fmt.Fprintln(w, "		return err")
	fmt.Fprintln(w, "	}")
	fmt.Fprintln(w, "")

	// todo: use multierror
	sa.IterateFields(func(fa *accessor.FieldAccessor) error {
		if !fa.Exported() {
			return nil
		}
		switch fa.IsRequired() {
		case true:
			fmt.Fprintf(w, "	if p.%s == nil {\n", fa.Name())
			fmt.Fprintf(w, "		return errors.New(\"%s is required\")\n", fa.GuessJSONFieldName(fa.Name()))
			fmt.Fprintf(w, "	}\n")
			fmt.Fprintf(w, "	x.%s = *p.%s\n", fa.Name(), fa.Name())
		case false:
			fmt.Fprintf(w, "	if p.%s != nil {\n", fa.Name())
			fmt.Fprintf(w, "		x.%s = *p.%s\n", fa.Name(), fa.Name())
			fmt.Fprintf(w, "	}\n")
		}
		return nil
	})
	fmt.Fprintln(w, "	return nil")
	return nil
}

// NameTo :
func NameTo(pkg *types.Package, f *ast.File) types.Qualifier {
	return func(other *types.Package) string {
		if pkg == other {
			return "" // same package; unqualified
		}
		// todo: cache
		for _, is := range f.Imports {
			if is.Path.Value[1:len(is.Path.Value)-1] == other.Path() {
				if is.Name != nil {
					return is.Name.String()
				}
				return other.Name()
			}
		}
		return other.Name() // todo: add import
	}
}
