package formatchecktask

import (
	"bytes"
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
	var b bytes.Buffer
	isEmpty := true
	o := writerutil.LeveledOutput{W: &b}

	o.Printf("// FormatCheck : (generated from %s)\n", ob.Type().String())
	o.WithBlock(fmt.Sprintf("func (x %s) FormatCheck() error", types.TypeString(types.NewPointer(ob.Type()), qf)), func() {
		o.Println("var merr *multierror.Error")
		o.Newline()

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
					isEmpty = false
					o.Printf("merr = multierror.Append(merr, errors.WithMessage(err, %q))\n", fa.Object.Name())
				})
			}
			if t, ok := fa.Object.Type().Underlying().(hasElem); ok {
				if formatcheckable.IsFormatCheckable(t.Elem()) {
					if _, ok := t.(*types.Pointer); !ok {
						o.WithBlock(fmt.Sprintf("for i, sub := range x.%s", fa.Object.Name()), func() {
							o.WithBlock("if err := sub.FormatCheck(); err != nil", func() {
								isEmpty = false
								o.Printf("merr = multierror.Append(merr, errors.WithMessage(err, fmt.Sprintf(\"%s[%%v]\", i)))\n", fa.Object.Name())
							})
						})
					}
				}
			}

			for _, c := range candidates {
				if v, ok := fa.Tag.Lookup(c.Name); ok {
					c := c // copied
					c.Value = v
					isEmpty = false
					c.Callback(&o, "x", fa, c.Value)
				}
			}
			return nil
		})

		o.Println("return merr.ErrorOrNil()")
	})

	if !isEmpty {
		io.Copy(w, &b)
		return nil
	}

	{
		o := writerutil.LeveledOutput{W: w}
		o.Printf("// FormatCheck : (generated from %s)\n", ob.Type().String())
		o.WithBlock(fmt.Sprintf("func (x %s) FormatCheck() error", types.TypeString(types.NewPointer(ob.Type()), qf)), func() {
			o.Println("return nil")
		})
	}
	return nil
}
