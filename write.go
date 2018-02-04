package strangejson

import (
	"fmt"
	"go/types"
	"io"
	"reflect"
)

// WritePackage :
func WritePackage(pkg *types.Package, w io.Writer) error {
	s := pkg.Scope()
	for _, name := range s.Names() {
		ob := s.Lookup(name)

		if internal, ok := ob.Type().Underlying().(*types.Struct); ok {
			fmt.Printf("%s\n", name)
			for i := 0; i < internal.NumFields(); i++ {
				dbname, found := reflect.StructTag(internal.Tag(i)).Lookup("db")
				field := internal.Field(i)
				fmt.Printf("	%v (exported=%t, dbname=%s, found=%t)\n", field, field.Exported(), dbname, found)
			}
			fmt.Println("")
		}
	}

	fmt.Fprintln(w, pkg.Name())
	return nil
}
