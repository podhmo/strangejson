package accessor

import (
	"go/types"
)

// Accessor :
type Accessor struct {
	Pkg *types.Package
}

// Lookup :
func (a *Accessor) Lookup(name string) types.Object {
	return a.Pkg.Scope().Lookup(name)
}

// IterateStructs :
func (a *Accessor) IterateStructs(
	fn func(sa *StructAccessor) error,
) error {
	scope := a.Pkg.Scope()

	for _, name := range scope.Names() {
		ob := scope.Lookup(name)

		if typename, ok := ob.(*types.TypeName); ok {
			if typename.IsAlias() {
				continue
			}
		}

		if underlying, ok := ob.Type().Underlying().(*types.Struct); ok {
			sa := &StructAccessor{
				Object:     ob,
				Underlying: underlying,
				fields:     map[string]*FieldAccessor{},
			}
			if err := fn(sa); err != nil {
				return err
			}
		}
	}
	return nil
}
