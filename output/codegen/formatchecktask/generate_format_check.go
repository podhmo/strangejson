package formatchecktask

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/strangejson/formatcheck"
	"github.com/podhmo/strangejson/output/codegen/accessor"
)

func generateFormatCheck(w io.Writer, f *ast.File, a *accessor.Accessor, sa *accessor.StructAccessor, qf types.Qualifier, formatcheckable *FormatCheckable) error {
	type hasElem interface {
		Elem() types.Type
	}

	ob := sa.Object

	fmt.Fprintf(w, "// FormatCheck : (generated from %s)\n", ob.Type().String())
	fmt.Fprintf(w, "func (x %s) FormatCheck() error {\n", types.TypeString(types.NewPointer(ob.Type()), qf))
	defer fmt.Fprintln(w, "}")

	// todo: use multierror
	candidates := []formatcheck.Check{formatcheck.MaxLength, formatcheck.MinLength}
	sa.IterateFields(func(fa *accessor.FieldAccessor) error {
		// todo: more sophisticated
		if _, ok := fa.Object.Type().(*types.Pointer); ok {
			fmt.Fprintf(w, "	if x.%s != nil {\n", fa.Object.Name())
			if formatcheckable.IsFormatCheckable(fa.Object.Type()) {
				fmt.Fprintf(w, "		if err := x.%s.FormatCheck(); err != nil {\n", fa.Object.Name())
				fmt.Fprintln(w, "			return err")
				fmt.Fprintln(w, "		}")
			}

			if t, ok := fa.Object.Type().Underlying().(hasElem); ok {
				if formatcheckable.IsFormatCheckable(t.Elem()) {
					fmt.Fprintf(w, "		for _, sub := range x.%s {\n", fa.Object.Name())
					fmt.Fprintln(w, "			if err := sub.FormatCheck(); err != nil {")
					fmt.Fprintln(w, "				return err")
					fmt.Fprintln(w, "			}")
					fmt.Fprintln(w, "		}")
				}
			}
			for _, c := range candidates {
				if v, ok := fa.Tag.Lookup(c.Name); ok {
					c := c // copied
					c.Value = v
					c.Callback(w, "x", fa, c.Value)
				}
			}
			fmt.Fprintln(w, "	}")
		} else {
			if formatcheckable.IsFormatCheckable(fa.Object.Type()) {
				fmt.Fprintf(w, "	if err := x.%s.FormatCheck(); err != nil {\n", fa.Object.Name())
				fmt.Fprintln(w, "		return err")
				fmt.Fprintln(w, "	}")
			}

			if t, ok := fa.Object.Type().Underlying().(hasElem); ok {
				if formatcheckable.IsFormatCheckable(t.Elem()) {
					fmt.Fprintf(w, "	for _, sub := range x.%s {\n", fa.Object.Name())
					fmt.Fprintln(w, "		if err := sub.FormatCheck(); err != nil {")
					fmt.Fprintln(w, "			return err")
					fmt.Fprintln(w, "		}")
					fmt.Fprintln(w, "	}")
				}
			}

			for _, c := range candidates {
				if v, ok := fa.Tag.Lookup(c.Name); ok {
					c := c // copied
					c.Value = v
					c.Callback(w, "x", fa, c.Value)
				}
			}
		}
		return nil
	})

	fmt.Fprintln(w, "	return nil")
	return nil
}
