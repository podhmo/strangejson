package unmarshaljsontask

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/strangejson/output/codegen/accessor"
	"github.com/podhmo/strangejson/output/codegen/writerutil"
)

func generateUnmarshalJSON(w io.Writer, f *ast.File, a *accessor.Accessor, sa *accessor.StructAccessor, qf types.Qualifier, sameOb types.Object) error {
	ob := sa.Object
	o := writerutil.LeveledOutput{W: w}

	o.Printf("// UnmarshalJSON : (generated from %s)\n", ob.Type().String())
	o.WithBlock(fmt.Sprintf("func (x %s) UnmarshalJSON(b []byte) error", types.TypeString(types.NewPointer(ob.Type()), qf)), func() {
		if sameOb != nil {
			o.Printf("return (%s)(x).UnmarshalJSON(b)\n", types.TypeString(types.NewPointer(sameOb.Type()), qf))
			return
		}

		// internal struct, all fields are pointer
		o.WithBlock("type internal struct", func() {
			sa.IterateFields(func(fa *accessor.FieldAccessor) error {
				if !fa.Exported() {
					return nil
				}
				typ := types.TypeString(types.NewPointer(fa.Object.Type()), qf)
				switch fa.Tag {
				case "":
					o.Printf("%s %s\n", fa.Name(), typ)
				default:
					o.Printf("%s %s `%s`\n", fa.Name(), typ, fa.Tag)
				}
				return nil
			})
		})

		o.Newline()
		o.Println("var p internal")

		o.WithBlock("if err := json.Unmarshal(b, &p); err != nil", func() {
			o.Println("return err")
		})
		o.Newline()

		// todo: use multierror
		sa.IterateFields(func(fa *accessor.FieldAccessor) error {
			if !fa.Exported() {
				return nil
			}
			switch fa.IsRequired() {
			case true:
				o.WithBlock(fmt.Sprintf("if p.%s == nil", fa.Name()), func() {
					o.Printf("return errors.New(\"%s is required\")\n", fa.GuessJSONFieldName(fa.Name()))
				})
				o.Printf("x.%s = *p.%s\n", fa.Name(), fa.Name())
			case false:
				o.WithBlock(fmt.Sprintf("if p.%s != nil", fa.Name()), func() {
					o.Printf("x.%s = *p.%s\n", fa.Name(), fa.Name())
				})
			}
			return nil
		})
		o.Println("return x.FormatCheck()")
	})
	return nil
}
