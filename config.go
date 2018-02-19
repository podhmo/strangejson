package strangejson

import (
	"go/build"
	"go/parser"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

// Config :
type Config struct {
	*loader.Config
}

// NewConfig :
func NewConfig(options ...func(c *Config)) *Config {
	c := &Config{
		Config: &loader.Config{
			Build:       &build.Default,
			ParserMode:  parser.ParseComments,
			AllowErrors: false,
			TypeCheckFuncBodies: func(path string) bool {
				return false // skip all type-check
			},
		},
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

// WithBuildContext :
func WithBuildContext(ctxt *build.Context) func(c *Config) {
	return func(c *Config) {
		c.Config.Build = ctxt
	}
}

// WithAllowErrors :
func WithAllowErrors(status bool) func(c *Config) {
	return func(c *Config) {
		c.Config.AllowErrors = status
	}
}

// ImportPkg :
func (conf *Config) ImportPkg(pkgpath string) []string {
	if strings.HasSuffix(pkgpath, "/*") {
		return conf.importPkgDir(strings.TrimSuffix(pkgpath, "/*"))
	}
	conf.Import(pkgpath)
	return []string{pkgpath}
}

func (conf *Config) importPkgDir(rootPkg string) []string {
	var subPkgs []string
	for _, dir := range conf.Build.SrcDirs() {
		fs, err := ioutil.ReadDir(filepath.Join(dir, rootPkg))
		_ = err
		for _, f := range fs {
			if f.IsDir() {
				pkgpath := filepath.Join(rootPkg, f.Name())
				subPkgs = append(subPkgs, pkgpath)
				conf.Import(pkgpath)
			}
		}
	}
	return subPkgs
}
