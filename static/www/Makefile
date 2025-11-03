LEAFLET_VERSION=1.9.4
ESRI_LEAFLET_VERSION=3.0.19
ESRI_LEAFLET_VECTOR_VERSION=4.3.2

PROTOMAPS_LEAFLET_VERSION=4.0.1
PROTOMAPS_LEAFLET_VERSION=5.1.0

leaflet:
	curl -s -o css/leaflet-$(LEAFLET_VERSION).css https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css
	curl -s -o javascript/leaflet-$(LEAFLET_VERSION).js https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js

esri-leaflet:
	curl -s -o javascript/esri-leaflet-$(ESRI_LEAFLET_VERSION).js https://unpkg.com/esri-leaflet@$(ESRI_LEAFLET_VERSION)/dist/esri-leaflet.js
	curl -s -o javascript/esri-leaflet-vector-$(ESRI_LEAFLET_VECTOR_VERSION).js https://unpkg.com/esri-leaflet-vector@$(ESRI_LEAFLET_VECTOR_VERSION)/dist/esri-leaflet-vector.js

protomaps:
	curl -s -o javascript/protomaps-leaflet-$(PROTOMAPS_LEAFLET_VERSION).js https://unpkg.com/protomaps-leaflet@$(PROTOMAPS_LEAFLET_VERSION)/dist/protomaps-leaflet.js

wasmjs:
	curl -s -o javascript/sfomuseum.golang.wasm.bundle.js https://raw.githubusercontent.com/sfomuseum/js-sfomuseum-golang-wasm/refs/heads/main/dist/sfomuseum.golang.wasm.bundle.js
