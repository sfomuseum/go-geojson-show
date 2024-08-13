GOMOD=vendor

cli:
	go build -mod $(GOMOD) -ldflags="-s -w"  -o bin/example cmd/example/main.go

debug:
	go run -mod $(GOMOD) cmd/example/main.go -javascript-at-eof -rollup-assets

protomaps-js:
	curl -s -L -o static/javascript/protomaps.js https://unpkg.com/protomaps@latest/dist/protomaps.js 
	curl -s -L -o static/javascript/protomaps.min.js https://unpkg.com/protomaps@latest/dist/protomaps.min.js 
