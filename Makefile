test:
	(cd cmd; go install -v ./...)
	(cd output; go test ./...)
	$(MAKE) -C examples

gen:
	(cd cmd; go install -v ./...)
	strangejson --pkg github.com/podhmo/strangejson/examples/simple00
	strangejson --pkg github.com/podhmo/strangejson/examples/depends01
	strangejson --pkg github.com/podhmo/strangejson/examples/pointer02
	strangejson --pkg github.com/podhmo/strangejson/examples/manytypes03
	strangejson --pkg github.com/podhmo/strangejson/examples/manypackages04/*

clean:
	find examples -name "*_gen.go" | xargs rm -vf
