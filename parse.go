package strangejson

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"golang.org/x/tools/go/ast/astutil"
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
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Type        Type       `json:"type"`
	Properties  []Property `json:"properties,omitempty"`
	Required    []string   `json:"required,omitempty"`

	XGoName string   `json:"x-goname"`
	Depends []Schema `json:"-"`
}

// Property :
type Property struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        Type   `json:"type"`
	Required    bool   `json:"required"`
	XGoName     string `json:"x-goname"`
}

// FindFileByPos :
func FindFileByPos(files []*ast.File, pos token.Pos) *ast.File {
	var found *ast.File
	for _, f := range files {
		if pos >= f.Pos() {
			found = f
		} else {
			return found
		}
	}
	return found
}

// FindDocStringByPos :
func FindDocStringByPos(files []*ast.File, pos token.Pos) *ast.CommentGroup {
	file := FindFileByPos(files, pos)
	if file == nil {
		return nil
	}

	nodes, _ := astutil.PathEnclosingInterval(file, pos, pos)
	if len(nodes) <= 0 {
		return nil
	}

	switch t := nodes[0].(type) {
	case *ast.GenDecl:
		return t.Doc
	case *ast.Ident:
		if t.Obj == nil {
			return nil
		}
		switch x := t.Obj.Decl.(type) {
		case *ast.Field:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		case *ast.ImportSpec:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		case *ast.ValueSpec:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		case *ast.TypeSpec:
			pp.Println(x)
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		default:
			log.Printf("default2: %#v\n", x)
			return nil
		}
	default:
		log.Printf("default: %#v\n", t)
		return nil
	}
}

// ParsePackageInfo :
func ParsePackageInfo(info *loader.PackageInfo, findDescription bool) ([]Schema, error) {
	sort.SliceStable(info.Files, func(i int, j int) bool {
		return info.Files[i].Pos() <= info.Files[j].Pos()
	})
	scope := info.Pkg.Scope()

	// todo: parse depends

	var r []Schema
	for _, name := range scope.Names() {
		ob := scope.Lookup(name)
		internal, ok := ob.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		s := Schema{
			Name:       name,
			Type:       TypeObject,
			Properties: []Property{},
			Required:   []string{},
			XGoName:    name,
		}

		if findDescription {
			description := FindDocStringByPos(info.Files, ob.Pos()-1)
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

			prop := Property{
				Name:     fieldname,
				Type:     ParseType(field),
				Required: required,
				XGoName:  field.Name(),
			}

			// if findDescription {
			// 	description := FindDocStringByPos(info.Files, field.Pos())
			// 	if description != nil {
			// 		prop.Description = strings.Trim(strings.TrimPrefix(description.Text(), prop.XGoName), " :\n")
			// 	}
			// }
			s.Properties = append(s.Properties, prop)
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
		case types.String, types.Byte: // byte?
			return TypeString
		case types.Float32, types.Float64:
			return TypeNumber
		case types.Int, types.Int16, types.Int32, types.Int64, types.Int8, types.Uint, types.Uint16, types.Uint32, types.Uint64:
			return TypeInteger
		case types.Bool:
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
