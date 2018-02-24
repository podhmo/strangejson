package formatchecktask

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/strangejson/formatcheck"
	"github.com/podhmo/strangejson/output/codegen/accessor"
	"github.com/podhmo/strangejson/output/codegen/writerutil"
)

func generateFormatCheck(w io.Writer, f *ast.File, a *accessor.Accessor, sa *accessor.StructAccessor, qf types.Qualifier, formatcheckable *FormatCheckable) error {
	type hasElem interface {
		Elem() types.Type
	}

	ob := sa.Object
	o := writerutil.LeveledOutput{W: w}

	o.Printf("// FormatCheck : (generated from %s)\n", ob.Type().String())
	o.WithBlock(fmt.Sprintf("func (x %s) FormatCheck() error", types.TypeString(types.NewPointer(ob.Type()), qf)), func() {

		// todo: use multierror
		candidates := []formatcheck.Check{formatcheck.MaxLength, formatcheck.MinLength}
		sa.IterateFields(func(fa *accessor.FieldAccessor) error {
			if _, ok := fa.Object.Type().(*types.Pointer); ok {
				o.Printf("if x.%s != nil {\n", fa.Object.Name())
				o.Indent()
				defer func() {
					o.UnIndent()
					o.Println("}")
				}()
			}

			if formatcheckable.IsFormatCheckable(fa.Object.Type()) {
				o.WithBlock(fmt.Sprintf("if err := x.%s.FormatCheck(); err != nil", fa.Object.Name()), func() {
					o.Println("	return err")
				})
			}
			if t, ok := fa.Object.Type().Underlying().(hasElem); ok {
				if formatcheckable.IsFormatCheckable(t.Elem()) {
					if _, ok := t.(*types.Pointer); !ok {
						o.WithBlock(fmt.Sprintf("for _, sub := range x.%s", fa.Object.Name()), func() {
							o.WithBlock("if err := sub.FormatCheck(); err != nil", func() {
								o.Println("return err")
							})
						})
					}
				}
			}

			for _, c := range candidates {
				if v, ok := fa.Tag.Lookup(c.Name); ok {
					c := c // copied
					c.Value = v
					c.Callback(&o, "x", fa, c.Value)
				}
			}
			return nil
		})

		o.Println("return nil")
	})
	return nil
}
