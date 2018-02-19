package codegen

import (
	"testing"

	"github.com/podhmo/strangejson"
	"github.com/podhmo/strangejson/buildcontext"
	"github.com/stretchr/testify/require"
)

func TestCodegen(t *testing.T) {
	pkgs := map[string]map[string]string{
		"github.com/podhmo/sandbox/model": {
			"person.go": `
package model
type Person struct {
	Name string
	Age int
}
`,
		},
	}
	build := buildcontext.FakeContext(pkgs)
	conf := strangejson.NewConfig(strangejson.WithBuildContext(build.Context))

	pkgpaths := conf.ImportPkg("github.com/podhmo/sandbox/model")
	prog, err := conf.Load()

	require.NoError(t, err)
	cmd := New(build, prog)
	require.NoError(t, cmd.Run(pkgpaths))

	t.Log(pkgs["github.com/podhmo/sandbox/model"]["person_gen.go"])
}
