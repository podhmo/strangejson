package main

import (
	"fmt"
	"log"
	"os"

	"github.com/podhmo/strangejson"
	"github.com/podhmo/strangejson/buildcontext"
	"github.com/podhmo/strangejson/output/codegen"
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
	build := buildcontext.Default()
	conf := strangejson.NewConfig(strangejson.WithBuildContext(build.Context))
	pkgpaths := conf.ImportPkg(opt.Pkg)
	prog, err := conf.Load()
	if err != nil {
		return err
	}
	// cmd := swaggergen.New(prog)
	cmd := codegen.New(build, prog)
	return cmd.Run(pkgpaths)
}
