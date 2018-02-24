package formatcheck

import (
	"fmt"

	"github.com/podhmo/strangejson/output/codegen/accessor"
	"github.com/podhmo/strangejson/output/codegen/writerutil"
)

// Check : todo: rename
type Check struct {
	Name     string
	Value    string
	Callback func(o *writerutil.LeveledOutput, name string, fa *accessor.FieldAccessor, val string) error
}

// MaxLength :
var MaxLength = Check{
	Name: "maxLength",
	Callback: func(o *writerutil.LeveledOutput, name string, fa *accessor.FieldAccessor, val string) error {
		o.WithBlock(fmt.Sprintf("if len(%s.%s) > %s", name, fa.Name(), val), func() {
			o.Println("return errors.New(\"max\")")
		})
		return nil
	},
}

// MinLength :
var MinLength = Check{
	Name: "minLength",
	Callback: func(o *writerutil.LeveledOutput, name string, fa *accessor.FieldAccessor, val string) error {
		o.WithBlock(fmt.Sprintf("if len(%s.%s) < %s", name, fa.Name(), val), func() {
			o.Println("return errors.New(\"min\")")
		})
		return nil
	},
}
