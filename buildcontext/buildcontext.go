package buildcontext

import (
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/buildutil"
)

// Context : (from *build.Context, not context.Context)
type Context struct {
	*build.Context
	WriteFile func(filename string, data []byte, perm os.FileMode) error
}

// WriteFile :
func WriteFile(build *Context, filename string, data []byte, perm os.FileMode) error {
	if build.WriteFile != nil {
		return build.WriteFile(filename, data, perm)
	}
	return ioutil.WriteFile(filename, data, perm)
}

// Default :
func Default() *Context {
	return &Context{Context: &build.Default}
}

// FakeContext :
func FakeContext(pkgs map[string]map[string]string) *Context {
	// copy from golang.org/x/tools/go/buildutil/facecontext.go
	clean := func(filename string) string {
		f := path.Clean(filepath.ToSlash(filename))
		// Removing "/go/src" while respecting segment
		// boundaries has this unfortunate corner case:
		if f == "/go/src" {
			return ""
		}
		return strings.TrimPrefix(f, "/go/src/")
	}

	ctx := &Context{Context: buildutil.FakeContext(pkgs)}
	ctx.WriteFile = func(filename string, data []byte, perm os.FileMode) error {
		filename = clean(filename)
		dir, base := path.Split(filename)

		dir = path.Clean(dir)
		if _, ok := pkgs[dir]; !ok {
			pkgs[dir] = map[string]string{}
		}

		pkgs[dir][base] = string(data)
		return nil
	}
	return ctx
}
