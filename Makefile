test:
	(cd cmd; go install -v ./...)
	(cd output/codegen; go test -v)

gen:
	(cd cmd; go install -v ./...)
	# strangejson --pkg github.com/podhmo/strangejson/examples/simple00
	strangejson --pkg github.com/podhmo/strangejson/examples/depends01

clean:
	find examples -name "*_gen.go" | xargs rm -vf
