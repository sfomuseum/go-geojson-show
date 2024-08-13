GOMOD=vendor

example:
	go run -mod $(GOMOD) cmd/example/main.go \
	-enable-hash \
	-enable-fullscreen \
	-enable-draw \
	-rollup-assets \
	-javascript-at-eof \
	-tile-url 'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
