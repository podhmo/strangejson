package accessor

import (
	"go/types"
	"reflect"
	"strconv"
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
func (a *Accessor) IterateStructs(fn func(sa *StructAccessor) error) error {
	scope := a.Pkg.Scope()
	for _, name := range scope.Names() {
		ob := scope.Lookup(name)
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

// StructAccessor :
type StructAccessor struct {
	Object     types.Object
	Underlying *types.Struct
	fields     map[string]*FieldAccessor // lazily allocated
}

// Exported :
func (sa *StructAccessor) Exported() bool {
	return sa.Object.Exported()
}

// Name :
func (sa *StructAccessor) Name() string {
	return sa.Object.Name()
}

// LookupFieldOrMethod :
func (sa *StructAccessor) LookupFieldOrMethod(name string) *FieldAccessor {
	if fa, ok := sa.fields[name]; ok {
		return fa
	}
	ob := sa.Object
	fieldOrMethod, _, _ := types.LookupFieldOrMethod(ob.Type(), true, ob.Pkg(), name)
	if fieldOrMethod == nil {
		return nil
	}

	fa := &FieldAccessor{Object: fieldOrMethod}
	for i := 0; i < sa.Underlying.NumFields(); i++ {
		if sa.Underlying.Field(i) == fieldOrMethod {
			fa.Tag = reflect.StructTag(sa.Underlying.Tag(i))
			break
		}
	}
	sa.fields[name] = fa
	return fa
}

// IterateFields :
func (sa *StructAccessor) IterateFields(fn func(sa *FieldAccessor) error) error {
	for i := 0; i < sa.Underlying.NumFields(); i++ {
		field := sa.Underlying.Field(i)
		name := field.Name()
		fa, ok := sa.fields[name]
		if !ok {
			fa = &FieldAccessor{Object: field, Tag: reflect.StructTag(sa.Underlying.Tag(i))}
			sa.fields[name] = fa
		}
		if err := fn(fa); err != nil {
			return err
		}
	}
	return nil
}

// FieldAccessor :
type FieldAccessor struct {
	Object types.Object
	Tag    reflect.StructTag
}

// Exported :
func (fa *FieldAccessor) Exported() bool {
	return fa.Object.Exported()
}

// Name :
func (fa *FieldAccessor) Name() string {
	return fa.Object.Name()
}

// IsRequired :
func (fa *FieldAccessor) IsRequired() bool {
	requiredStr, ok := fa.Tag.Lookup("required")
	if !ok {
		requiredStr = "true"
	}
	required, err := strconv.ParseBool(requiredStr)
	if err != nil {
		required = true
	}
	return required
}

// GuessJSONFieldName :
func (fa *FieldAccessor) GuessJSONFieldName(defaultStr string) string {
	// todo: support omitempty
	fieldname, ok := fa.Tag.Lookup("json")
	if !ok {
		return defaultStr
	}
	return fieldname
}
