test:
	(cd cmd; go install -v ./...)
	(cd output; go test ./...)
	$(MAKE) -C _examples

gen:
	(cd cmd; go install -v ./...)
	strangejson --pkg github.com/podhmo/strangejson/_examples/simple00
	strangejson --pkg github.com/podhmo/strangejson/_examples/depends01
	strangejson --pkg github.com/podhmo/strangejson/_examples/pointer02
	strangejson --pkg github.com/podhmo/strangejson/_examples/manytypes03
	strangejson --pkg github.com/podhmo/strangejson/_examples/manypackages04/*
	strangejson --pkg github.com/podhmo/strangejson/_examples/formatcheck05

clean:
	find examples -name "*_gen.go" | xargs rm -vf
