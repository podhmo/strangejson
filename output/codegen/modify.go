package codegen

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/podhmo/astknife/action"
	"github.com/podhmo/astknife/lookup"
	"github.com/podhmo/strangejson/output"
)

// Modifier :
type Modifier struct {
	fset   *token.FileSet
	f      *ast.File
	lookup *lookup.Lookup
}

func (m *Modifier) modifyCode(filename string, gens []output.Gen, body []byte) ([]byte, error) {
	newf, err := parser.ParseFile(m.fset, filename, body, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if err := m.modifyAST(newf, gens); err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err := printer.Fprint(&b, m.fset, m.f); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (m *Modifier) modifyAST(newf *ast.File, gens []output.Gen) error {
	for _, gen := range gens {
		replacement := m.lookup.With(newf).Lookup(gen.Name) // xxx
		if _, err := action.AppendOrReplace(m.lookup, m.f, replacement); err != nil {
			if !action.IsNoEffect(err) {
				return err
			}
		}
	}
	return nil
}
