GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
CWD=$(shell pwd)

INITIAL_VIEW=-122.384292,37.621131,13

example:
	go run cmd/example/main.go \
		-initial-view '$(INITIAL_VIEW)'

example-protomaps:
	go run cmd/example/main.go \
		-initial-view '$(INITIAL_VIEW)' \
		-map-provider protomaps \
		-map-tile-uri 'file://$(CWD)/fixtures/sfo.pmtiles'
