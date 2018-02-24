package unmarshaljsontask

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/strangejson/output/codegen/accessor"
)

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
	fmt.Fprintln(w, "	return x.FormatCheck()")
	return nil
}
