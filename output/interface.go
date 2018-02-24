package output

import (
	"go/ast"
	"go/types"
	"io"

	"github.com/podhmo/astknife/bypos"
)

// Command :
type Command interface {
	Run(pkgpaths []string) error
}

// Task :
type Task interface {
	Prepare(pkg *types.Package, files bypos.Sorted) error

	Do(pkg *types.Package, files bypos.Sorted) ([]Gen, error)
}

// Gen :
type Gen struct {
	Name     string
	Object   types.Object
	File     *ast.File
	Generate func(w io.Writer, f *ast.File) error // xxx tentative
}
