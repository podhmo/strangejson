gen:
	go install -v ../../cmd/...
	strangejson --pkg github.com/podhmo/strangejson/_examples/manypackages04/*

clean:
	find . -name "*_gen.go" | xargs rm -vf
