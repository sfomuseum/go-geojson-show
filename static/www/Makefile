BOOTSTRAP_VERSION=5.3.8
LEAFLET_VERSION=1.9.4
PROTOMAPS_LEAFLET_VERSION=4.0.1

bootstrap:
	curl -s -o css/bootstrap.min.css https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/css/bootstrap.min.css
	curl -s -o javascript/bootstrap.bundle.min.js https://cdn.jsdelivr.net/npm/bootstrap@$(BOOTSTRAP_VERSION)/dist/js/bootstrap.bundle.min.js

leaflet:
	curl -s -o css/leaflet.css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet.js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

protomaps:
	curl -s -o javascript/protomaps-leaflet.js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION)/dist/protomaps-leaflet.js

wasmjs:
	curl -s -o javascript/sfomuseum.golang.wasm.bundle.js https://raw.githubusercontent.com/sfomuseum/js-sfomuseum-golang-wasm/refs/heads/main/dist/sfomuseum.golang.wasm.bundle.js
