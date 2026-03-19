GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
CWD=$(shell pwd)

INITIAL_VIEW=-122.384292,37.621131,13

example:
	go run cmd/example/main.go \
		-leaflet-pane hello=100 \
		-leaflet-pane world=200 \
		-initial-view '$(INITIAL_VIEW)'

example-pm:
	go run cmd/example/main.go \
		-initial-view '$(INITIAL_VIEW)' \
		-map-provider protomaps \
		-protomaps-max-data-zoom 14 \
		-map-tile-uri 'file://$(CWD)/fixtures/sfo.pmtiles'

example-pm-paint:
	go run cmd/example/main.go \
		-initial-view '$(INITIAL_VIEW)' \
		-map-provider protomaps-paint \
		-protomaps-max-data-zoom 14 \
		-map-tile-uri 'file://$(CWD)/fixtures/sfo.pmtiles'

example-pm-raster:
	go run cmd/example/main.go \
		-initial-view '$(INITIAL_VIEW)' \
		-map-provider protomaps-raster \
		-protomaps-max-data-zoom 14 \
		-map-tile-uri https://static.sfomuseum.org/aerial/1936.pmtiles

example-pm-ml:
	go run cmd/example/main.go \
		-initial-view '$(INITIAL_VIEW)' \
		-map-provider protomaps-ml \
		-protomaps-max-data-zoom 14 \
		-map-tile-uri 'file://$(CWD)/fixtures/sfo.pmtiles'
