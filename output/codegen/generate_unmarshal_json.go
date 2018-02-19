package codegen

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/strangejson/output/codegen/accessor"
)

// GenerateUnmarshalJSON :
func GenerateUnmarshalJSON(pkg *types.Package) ([]Gen, error) {
	a := &accessor.Accessor{Pkg: pkg}
	var results []Gen

	err := a.IterateStructs(func(sa *accessor.StructAccessor) error {
		if !sa.Exported() {
			return nil
		}
		results = append(results, Gen{
			Name:   "UnmarshalJSON",
			Object: sa.Object,
			Generate: func(w io.Writer, f *ast.File) error {
				qf := NameTo(pkg, f)
				return generateUnmarshalJSON(w, f, a, sa, qf)
			},
		})
		return nil
	})
	return results, err
}

func generateUnmarshalJSON(w io.Writer, f *ast.File, a *accessor.Accessor, sa *accessor.StructAccessor, qf types.Qualifier) error {
	unmarshalJSON := sa.LookupFieldOrMethod("UnmarshalJSON")
	if unmarshalJSON != nil {
		// TODO: update?
		return nil
	}

	ob := sa.Object
	typename := ob.Name()
	fmt.Fprintf(w, "// UnmarshalJSON : (generated from %s)\n", ob.Type().String())
	fmt.Fprintf(w, "func (x %s) UnmarshalJSON(b []byte) error {\n", typename)

	// internal struct, all fields are pointer
	fmt.Fprintf(w, "	type internal struct {\n")
	sa.IterateFields(func(fa *accessor.FieldAccessor) error {
		if !fa.IsRequired() {
			return nil
		}
		typ := types.TypeString(types.NewPointer(fa.Object.Type()), qf)
		if fa.Tag == "" {
			fmt.Fprintf(w, "		%s %s\n", fa.Name(), typ)
		} else {
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
		if !fa.IsRequired() {
			return nil
		}

		fmt.Fprintf(w, "	if p.%s == nil {\n", fa.Name())
		fmt.Fprintf(w, "		return errors.New(\"%s is required\")\n", fa.GuessJSONFieldName(fa.Name()))
		fmt.Fprintf(w, "	}\n")
		fmt.Fprintf(w, "	x.%s = *p.%s\n", fa.Name(), fa.Name())
		return nil
	})
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "	return nil")
	fmt.Fprintln(w, "}")
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
