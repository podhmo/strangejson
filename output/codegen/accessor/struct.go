package accessor

import (
	"go/types"
	"reflect"
)

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
