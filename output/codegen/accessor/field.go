package accessor

import (
	"go/types"
	"reflect"
	"strconv"
)

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
