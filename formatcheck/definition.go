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
			o.Printf("merr = multierror.Append(merr, errors.New(\"%s maxLength\"))\n", fa.GuessJSONFieldName(fa.Name()))
		})
		return nil
	},
}

// MinLength :
var MinLength = Check{
	Name: "minLength",
	Callback: func(o *writerutil.LeveledOutput, name string, fa *accessor.FieldAccessor, val string) error {
		o.WithBlock(fmt.Sprintf("if len(%s.%s) < %s", name, fa.Name(), val), func() {
			o.Printf("merr = multierror.Append(merr, errors.New(\"%s minLength\"))\n", fa.GuessJSONFieldName(fa.Name()))
		})
		return nil
	},
}

// Max :
var Max = Check{
	Name: "max",
	Callback: func(o *writerutil.LeveledOutput, name string, fa *accessor.FieldAccessor, val string) error {
		o.WithBlock(fmt.Sprintf("if %s.%s > %s", name, fa.Name(), val), func() {
			o.Printf("merr = multierror.Append(merr, errors.New(\"%s max\"))\n", fa.GuessJSONFieldName(fa.Name()))
		})
		return nil
	},
}

// Min :
var Min = Check{
	Name: "min",
	Callback: func(o *writerutil.LeveledOutput, name string, fa *accessor.FieldAccessor, val string) error {
		o.WithBlock(fmt.Sprintf("if %s.%s < %s", name, fa.Name(), val), func() {
			o.Printf("merr = multierror.Append(merr, errors.New(\"%s min\"))\n", fa.GuessJSONFieldName(fa.Name()))
		})
		return nil
	},
}
