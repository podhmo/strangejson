package strangejson

import (
	"go/types"
	"log"
	"reflect"
	"strconv"
)

// todo: namespace

// Type :
type Type string

// Type :
const (
	TypeObject  = Type("object")
	TypeArray   = Type("array")
	TypeBoolean = Type("boolean")
	TypeInteger = Type("integer")
	TypeNumber  = Type("number")
	TypeString  = Type("string")
)

// Schema :
type Schema struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Type        Type     `json:"type"`
	Properties  []Field  `json:"properties,omitempty"`
	Required    []string `json:"required,omitempty"`

	Depends []Schema `json:"-"`
}

// Field :
type Field struct {
	Name     string `json:"name"`
	Type     Type   `json:"type"`
	Required bool   `json:"required"`
	XGoName  string `json:"x-goname"`
}

// ParsePackage :
func ParsePackage(pkg *types.Package) ([]Schema, error) {
	scope := pkg.Scope()

	var r []Schema
	for _, name := range scope.Names() {
		ob := scope.Lookup(name)

		// todo: parse comment

		internal, ok := ob.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		s := Schema{
			Name:       name,
			Type:       TypeObject,
			Properties: []Field{},
			Required:   []string{},
		}

		for i := 0; i < internal.NumFields(); i++ {
			field := internal.Field(i)
			tag := reflect.StructTag(internal.Tag(i))

			// todo: omitempty
			fieldname, ok := tag.Lookup("json")
			if !ok {
				fieldname = field.Name()
			}

			requiredStr, ok := tag.Lookup("required")
			if !ok {
				requiredStr = "true"
			}
			required, err := strconv.ParseBool(requiredStr)
			if err != nil {
				required = true
			}

			s.Properties = append(s.Properties, Field{
				Name:     fieldname,
				Type:     ParseType(field),
				Required: required,
				XGoName:  field.Name(),
			})
		}
		for _, prop := range s.Properties {
			if prop.Required {
				s.Required = append(s.Required, prop.Name)
			}
		}
		r = append(r, s)
	}
	return r, nil
}

// ParseType :
func ParseType(ob types.Object) Type {
	switch t := ob.Type().Underlying().(type) {
	case *types.Struct:
		return TypeObject
	case *types.Array, *types.Slice:
		return TypeArray
	case *types.Basic:
		switch t.Kind() {
		case types.String, types.Byte: // byte?
			return TypeString
		case types.Float32, types.Float64:
			return TypeNumber
		case types.Int, types.Int16, types.Int32, types.Int64, types.Int8, types.Uint, types.Uint16, types.Uint32, types.Uint64:
			return TypeInteger
		default:
			log.Printf("unexpected type %#v (basic %#v)", ob.Type(), t)
			return TypeObject
		}
	default:
		log.Printf("unexpected type %#v", ob.Type())
		return TypeObject
	}
}
