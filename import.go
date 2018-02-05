package strangejson

import (
	"go/build"
	"go/parser"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

// NewConfig :
func NewConfig() *loader.Config {
	ctxt := build.Default
	conf := &loader.Config{
		Build:       &ctxt,
		ParserMode:  parser.ParseComments,
		AllowErrors: false, // xxx
	}
	return conf
}

// ImportPkg :
func ImportPkg(conf *loader.Config, pkgpath string) []string {
	if strings.HasSuffix(pkgpath, "/*") {
		return importPkgDir(conf, strings.TrimSuffix(pkgpath, "/*"))
	}
	conf.Import(pkgpath)
	return []string{pkgpath}
}

func importPkgDir(conf *loader.Config, rootPkg string) []string {
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
