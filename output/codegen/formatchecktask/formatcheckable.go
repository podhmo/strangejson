package formatchecktask

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/loader"
)

// FormatCheckable :
type FormatCheckable struct {
	iface              *types.Interface
	fakeImplementedMap map[string]struct{}
}

// NewFormatCheckable :
func NewFormatCheckable(prog *loader.Program) (*FormatCheckable, error) {
	pkg := prog.Package("github.com/podhmo/strangejson/formatcheck").Pkg

	ob := pkg.Scope().Lookup("FormatCheckable")
	if ob == nil {
		return nil, errors.New("FormatCheckable interface is not found (lookup)")
	}
	t, ok := ob.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, errors.New("FormatCheckable interface is not found (coerce)")
	}

	return &FormatCheckable{
		iface:              t,
		fakeImplementedMap: map[string]struct{}{},
	}, nil
}

// RegisterFake :
func (fc *FormatCheckable) RegisterFake(typ types.Type) {
	fc.fakeImplementedMap[typ.String()] = struct{}{}
	fc.fakeImplementedMap[types.NewPointer(typ).String()] = struct{}{}
}

// IsFormatCheckable :
func (fc *FormatCheckable) IsFormatCheckable(typ types.Type) bool {
	if _, ok := fc.fakeImplementedMap[typ.String()]; ok {
		return true
	}
	return types.Implements(typ, fc.iface) || types.Implements(types.NewPointer(typ), fc.iface)
}
