package main

import (
	"fmt"
	"go/build"
	"go/parser"
	"log"
	"os"

	"github.com/k0kubun/pp"
	"github.com/podhmo/strangejson"
	"github.com/podhmo/strangejson/finding"
	"golang.org/x/tools/go/loader"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type opt struct {
	Pkg string
}

func main() {
	var opt opt
	app := kingpin.New("strangejson", "strangejson")
	app.Flag("pkg", "package").Required().StringVar(&opt.Pkg)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(fmt.Sprintf("%v", err))
	}

	if err := run(&opt); err != nil {
		log.Fatal(err)
	}
}

func run(opt *opt) error {
	ctxt := build.Default

	conf := &loader.Config{
		Build:       &ctxt,
		ParserMode:  parser.ParseComments,
		AllowErrors: false, // xxx
	}

	pkgpaths := finding.ImportPkg(conf, opt.Pkg)
	prog, err := conf.Load()
	if err != nil {
		return err
	}

	for _, pkgpath := range pkgpaths {
		info := prog.Package(pkgpath)
		schemas, err := strangejson.ParsePackage(info.Pkg)
		if err != nil {
			return err
		}
		pp.Println(schemas)
	}
	return nil
}
