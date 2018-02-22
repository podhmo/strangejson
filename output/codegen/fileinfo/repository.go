package fileinfo

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Repository :
type Repository struct {
	Name string

	nameMap  map[*ast.File]string
	fileMap  map[string]*ast.File
	isNewMap map[*ast.File]struct{}
}

// NewRepository :
func NewRepository(fset *token.FileSet, pkg *types.Package, files []*ast.File) *Repository {
	filenameMap := make(map[*ast.File]string, len(files))
	fileMap := map[string]*ast.File{}
	for _, f := range files {
		filename := fset.Position(f.Pos()).Filename
		filenameMap[f] = filename
		fileMap[filename] = f
	}
	return &Repository{
		Name:     pkg.Name(),
		nameMap:  filenameMap,
		fileMap:  fileMap,
		isNewMap: map[*ast.File]struct{}{},
	}
}

// CreateFakeFile :
func (r *Repository) CreateFakeFile(name string) *ast.File {
	f := &ast.File{Name: &ast.Ident{Name: r.Name}}
	r.isNewMap[f] = struct{}{}
	return f
}

// AddFile :
func (r *Repository) AddFile(name string, f *ast.File) {
	r.nameMap[f] = name
	r.fileMap[name] = f
}

// NameOf :
func (r *Repository) NameOf(f *ast.File) (string, bool) {
	name, ok := r.nameMap[f]
	return name, ok
}

// FileOf :
func (r *Repository) FileOf(name string) (*ast.File, bool) {
	file, ok := r.fileMap[name]
	return file, ok
}

// IsNew :
func (r *Repository) IsNew(f *ast.File) bool {
	_, ok := r.isNewMap[f]
	return ok
}
