LEAFLET_VERSION=1.9.4
PROTOMAPS_LEAFLET_VERSION=4.0.1
PROTOMAPS_LEAFLET_VERSION=5.1.0

leaflet:
	curl -s -o css/leaflet-$(LEAFLET_VERSION).css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet-$(LEAFLET_VERSION).js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

protomaps:
	curl -s -o javascript/protomaps-leaflet-$(PROTOMAPS_LEAFLET_VERSION).js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION)/dist/protomaps-leaflet.js

wasmjs:
	curl -s -o javascript/sfomuseum.golang.wasm.bundle.js https://raw.githubusercontent.com/sfomuseum/js-sfomuseum-golang-wasm/refs/heads/main/dist/sfomuseum.golang.wasm.bundle.js
