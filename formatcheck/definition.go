package formatcheck

import (
	"fmt"
	"io"

	"github.com/podhmo/strangejson/output/codegen/accessor"
)

// Check : todo: rename
type Check struct {
	Name     string
	Value    string
	Callback func(w io.Writer, name string, fa *accessor.FieldAccessor, val string) error
}

// MaxLength :
var MaxLength = Check{
	Name: "maxLength",
	Callback: func(w io.Writer, name string, fa *accessor.FieldAccessor, val string) error {
		fmt.Fprintf(w, "	if len(%s.%s) > %s {\n", name, fa.Name(), val)
		fmt.Fprintf(w, "		return errors.New(\"max\")\n")
		fmt.Fprintf(w, "	}\n")
		return nil
	},
}

// MinLength :
var MinLength = Check{
	Name: "minLength",
	Callback: func(w io.Writer, name string, fa *accessor.FieldAccessor, val string) error {
		fmt.Fprintf(w, "	if len(%s.%s) < %s {\n", name, fa.Name(), val)
		fmt.Fprintf(w, "		return errors.New(\"min\")\n")
		fmt.Fprintf(w, "	}\n")
		return nil
	},
}
