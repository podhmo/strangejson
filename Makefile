gen:
	go install -v ./...
	# strangejson --pkg github.com/podhmo/strangejson/examples/simple00
	strangejson --pkg github.com/podhmo/strangejson/examples/depends01
