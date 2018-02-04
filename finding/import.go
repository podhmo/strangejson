package finding

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

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
