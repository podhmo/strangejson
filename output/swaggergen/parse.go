package swaggergen

import (
	"go/token"
	"go/types"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/podhmo/strangejson/astutil"
	"golang.org/x/tools/go/loader"
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
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Type        Type             `json:"type"`
	Properties  []SchemaProperty `json:"properties,omitempty"`
	Required    []string         `json:"required,omitempty"`

	XGoName string `json:"x-goname"`
	XGoType string `json:"x-gotype"`

	Enum []Enum    `json:"enum,omitempty"` //
	Pos  token.Pos `json:"-"`              // id
}

// Enum :
type Enum struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value"`

	XGoName string    `json:"x-goname"`
	Pos     token.Pos `json:"-"` // id
}

// SchemaProperty :
type SchemaProperty struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        Type   `json:"type"`
	Required    bool   `json:"required"`

	XGoName string `json:"x-goname"`
	XGoType string `json:"x-gotype"`

	Depends []token.Pos `json:"-"`
}

// ParsePackageInfo :
func ParsePackageInfo(info *loader.PackageInfo, findDescription bool) ([]*Schema, error) {
	sort.SliceStable(info.Files, func(i int, j int) bool {
		return info.Files[i].Pos() <= info.Files[j].Pos()
	})
	scope := info.Pkg.Scope()
	typeNameMap := map[string]*Schema{}

	// todo: parse dependencies

	var r []*Schema
	for _, name := range scope.Names() {
		ob := scope.Lookup(name)
		switch ob.Type().Underlying().(type) {
		case *types.Struct:
			r = append(r, ParseStruct(info, ob, findDescription))
		case *types.Slice:
			//
		case *types.Map:
			//
		case *types.Signature:
			//
		case *types.Interface:
			//
		case *types.Pointer: // ???
			//
		case *types.Basic:
			s, ok := typeNameMap[ob.Type().String()]
			if !ok {
				// xxx:
				s = &Schema{
					Name:    ob.Name(),
					Type:    ParseType(ob),
					XGoName: ob.Name(),
					XGoType: ob.Type().String(),
					Pos:     ob.Pos(),
				}

				if findDescription {
					description := astutil.FindDocStringByPos(info.Files, ob.Pos())
					if description != nil {
						s.Description = strings.Trim(strings.TrimPrefix(description.Text(), ob.Name()), " :\n")
					}
				}

				typeNameMap[ob.Type().String()] = s
				r = append(r, s)
			}
			switch v := ob.(type) {
			case *types.Const:
				enum := Enum{
					Name:    ob.Name(),
					XGoName: ob.Name(),
					Value:   v.Val().String(),
					Pos:     ob.Pos(),
				}
				if findDescription {
					description := astutil.FindDocStringByPos(info.Files, ob.Pos())
					if description != nil {
						enum.Description = strings.Trim(strings.TrimPrefix(description.Text(), ob.Name()), " :\n")
					}
				}
				s.Enum = append(s.Enum, enum)
			case *types.TypeName:
				// noop
			case *types.Var:
				// noop
			default:
				log.Printf(".. unexpected object %#v basic (underlying %#v)\n", ob, ob.Type().Underlying())
			}
		default:
			log.Printf("unexpected object %#v (underlying %#v)\n", ob, ob.Type().Underlying())
		}

	}
	return r, nil
}

// ParseStruct :
func ParseStruct(info *loader.PackageInfo, ob types.Object, findDescription bool) *Schema {
	name := ob.Name()
	internal := ob.Type().Underlying().(*types.Struct) // xxx

	s := &Schema{
		Name:       name,
		Type:       TypeObject,
		Properties: []SchemaProperty{},
		Required:   []string{},
		XGoName:    name,
		XGoType:    ob.Type().String(),
		Pos:        ob.Pos(),
	}

	if findDescription {
		description := astutil.FindDocStringByPos(info.Files, ob.Pos()-1)
		if description != nil {
			s.Description = strings.Trim(strings.TrimPrefix(description.Text(), s.XGoName), " :\n")
		}
	}

	for i := 0; i < internal.NumFields(); i++ {
		field := internal.Field(i)
		if !field.Exported() {
			continue
		}

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

		prop := SchemaProperty{
			Name:     fieldname,
			Type:     ParseType(field),
			Required: required,
			XGoName:  field.Name(),
			XGoType:  field.Type().String(),
		}

		if findDescription {
			description := astutil.FindDocStringByPos(info.Files, field.Pos())
			if description != nil {
				prop.Description = strings.Trim(strings.TrimPrefix(description.Text(), prop.XGoName), " :\n")
			}
		}
		s.Properties = append(s.Properties, prop)
	}
	for _, prop := range s.Properties {
		if prop.Required {
			s.Required = append(s.Required, prop.Name)
		}
	}
	return s
}

// ParseType :
func ParseType(ob types.Object) Type {
	return parseType(ob.Type())
}

func parseType(typ types.Type) Type {
	switch t := typ.(type) {
	case *types.Struct:
		return TypeObject
	case *types.Array, *types.Slice:
		return TypeArray
	case *types.Map:
		return TypeObject
	case *types.Basic:
		switch t.Kind() {
		case types.String, types.Byte, types.UntypedRune, types.UntypedString: // byte?
			return TypeString
		case types.Float32, types.Float64, types.UntypedFloat:
			return TypeNumber
		case types.Int, types.Int16, types.Int32, types.Int64, types.Int8, types.Uint, types.Uint16, types.Uint32, types.Uint64, types.UntypedInt:
			return TypeInteger
		case types.Bool, types.UntypedBool:
			return TypeBoolean
		default:
			log.Printf("unexpected type %#v (basic %#v)", typ, t)
			return TypeObject
		}
	case *types.Pointer:
		return parseType(t.Elem())
	case *types.Named:
		return parseType(t.Underlying())
	default:
		log.Printf("unexpected type %#v", typ)
		return TypeObject
	}
}
