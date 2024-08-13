package main

/*

> go run cmd/gh2b/main.go -mode geojson 9q8yy | go run cmd/show/main.go -
Features are viewable at http://localhost:8080

> go run cmd/gh2b/main.go -mode geojson 9q8yy 9q5ct | go run cmd/show/main.go -
Features are viewable at http://localhost:8080

go run cmd/gh2b/main.go -format geojson 9q8yp | go run cmd/show/main.go -map-provider protomaps -map-tile-uri file:///usr/local/sfomuseum/go-http-protomaps/cmd/example/sfo.pmtiles -

*/

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/paulmach/orb/geojson"
	"github.com/pkg/browser"
	"github.com/sfomuseum/go-http-protomaps"
	"github.com/tidwall/gjson"
	wasm_js "github.com/whosonfirst/go-whosonfirst-format-wasm/static/javascript"
	"github.com/whosonfirst/go-whosonfirst-format-wasm/static/wasm"
)

const leaflet_osm_tile_url = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
const protomaps_api_tile_url string = "https://api.protomaps.com/tiles/v3/{z}/{x}/{y}.mvt?key={key}"

//go:embed *.html
var html_FS embed.FS

//go:embed css/*.css
var css_FS embed.FS

//go:embed javascript/*.js javascript/*.map
var js_FS embed.FS

// mapConfig defines common configuration details for maps.
type mapConfig struct {
	// A valid map provider label.
	Provider string `json:"provider"`
	// A valid Leaflet tile layer URI.
	TileURL string `json:"tile_url"`
	// Optional Protomaps configuration details
	Protomaps *protomapsConfig `json:"protomaps,omitempty"`
}

// protomapsConfig defines configuration details for maps using Protomaps.
type protomapsConfig struct {
	// A valid Protomaps theme label
	Theme string `json:"theme"`
}

func main() {

	var port int

	var map_provider string
	var map_tile_uri string

	var protomaps_theme string

	flag.StringVar(&map_provider, "map-provider", "leaflet", "Valid options are: leaflet, protomaps")
	flag.StringVar(&map_tile_uri, "map-tile-uri", leaflet_osm_tile_url, "A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs.")

	flag.StringVar(&protomaps_theme, "protomaps-theme", "white", "A valid Protomaps theme label.")

	flag.IntVar(&port, "port", 0, "The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.")

	flag.Parse()

	ctx := context.Background()

	fc := geojson.NewFeatureCollection()

	append_features := func(r io.Reader) error {

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read body, %w", err)
		}

		type_rsp := gjson.GetBytes(body, "type")

		switch type_rsp.String() {
		case "Feature":

			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal Feature, %w", err)
			}

			fc.Append(f)

		case "FeatureCollection":

			other_fc, err := geojson.UnmarshalFeatureCollection(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal record as FeatureCollection, %w", err)
			}

			for _, f := range other_fc.Features {
				fc.Append(f)
			}

		default:
			return fmt.Errorf("Invalid type, %s", type_rsp.String())
		}

		return nil
	}

	uris := flag.Args()

	stdin := false

	if len(uris) == 1 && uris[0] == "-" {
		stdin = true
	}

	if stdin {

		err := append_features(os.Stdin)

		if err != nil {
			log.Fatalf("Failed to append features, %v", err)
		}

	} else {

		for _, path := range uris {

			r, err := os.Open(path)

			if err != nil {
				log.Fatalf("Failed to open %s for reading, %v", path, err)
			}

			defer r.Close()

			err = append_features(r)

			if err != nil {
				log.Fatalf("Failed to append features, %v", err)
			}
		}
	}

	//

	mux := http.NewServeMux()

	html_fs := http.FS(html_FS)
	js_fs := http.FS(js_FS)
	css_fs := http.FS(css_FS)

	mux.Handle("/css/", http.FileServer(css_fs))
	mux.Handle("/javascript/", http.FileServer(js_fs))
	mux.Handle("/", http.FileServer(html_fs))

	wasm_fs := http.FS(wasm.FS)
	wasm_handler := http.FileServer(wasm_fs)

	wasm_js_fs := http.FS(wasm_js.FS)
	wasm_js_handler := http.FileServer(wasm_js_fs)

	mux.Handle("/javascript/wasm/", http.StripPrefix("/javascript/wasm/", wasm_js_handler))
	mux.Handle("/wasm/", http.StripPrefix("/wasm/", wasm_handler))

	data_handler := dataHandler(fc)

	mux.Handle("/features.geojson", data_handler)

	//

	map_cfg := &mapConfig{
		Provider: map_provider,
		TileURL:  map_tile_uri,
	}

	if map_provider == "protomaps" {

		u, err := url.Parse(map_tile_uri)

		if err != nil {
			log.Fatalf("Failed to parse Protomaps tile URL, %w", err)
		}

		switch u.Scheme {
		case "file":

			mux_url, mux_handler, err := protomaps.FileHandlerFromPath(u.Path, "")

			if err != nil {
				log.Fatalf("Failed to determine absolute path for '%s', %v", map_tile_uri, err)
			}

			mux.Handle(mux_url, mux_handler)
			map_cfg.TileURL = mux_url

		case "api":
			key := u.Host
			map_cfg.TileURL = strings.Replace(protomaps_api_tile_url, "{key}", key, 1)
		}

		map_cfg.Protomaps = &protomapsConfig{
			Theme: protomaps_theme,
		}
	}

	map_cfg_handler := mapConfigHandler(map_cfg)

	mux.Handle("/map.json", map_cfg_handler)

	//

	if port == 0 {

		listener, err := net.Listen("tcp", "localhost:0")

		if err != nil {
			log.Fatalf("Failed to determine next available port, %v", err)
		}

		port = listener.Addr().(*net.TCPAddr).Port
		err = listener.Close()

		if err != nil {
			log.Fatalf("Failed to close listener used to derive port, %v", err)
		}
	}

	//

	addr := fmt.Sprintf("localhost:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	http_server := http.Server{
		Addr: addr,
	}

	http_server.Handler = mux

	done_ch := make(chan bool)
	err_ch := make(chan error)

	go func() {

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("Shutting server down")
		err := http_server.Shutdown(ctx)

		if err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
		}

		close(done_ch)
	}()

	go func() {

		err := http_server.ListenAndServe()

		if err != nil {
			err_ch <- fmt.Errorf("Failed to start server, %w", err)
		}
	}()

	server_ready := false

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-err_ch:
			log.Fatalf("Received error starting server, %v", err)
		case <-ticker.C:

			rsp, err := http.Head(url)

			if err != nil {
				slog.Warn("HEAD request failed", "url", url, "error", err)
			} else {

				defer rsp.Body.Close()

				if rsp.StatusCode != 200 {
					slog.Warn("HEAD request did not return expected status code", "url", url, "code", rsp.StatusCode)
				} else {
					slog.Debug("HEAD request succeeded", "url", url)
					server_ready = true
				}
			}
		}

		if server_ready {
			break
		}
	}

	err := browser.OpenURL(url)

	if err != nil {
		log.Fatalf("Failed to open URL %s, %v", url, err)
	}

	log.Printf("Features are viewable at %s\n", url)
	<-done_ch
}

func dataHandler(fc *geojson.FeatureCollection) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		enc_json, err := fc.MarshalJSON()

		if err != nil {
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write(enc_json)
		return
	}

	return http.HandlerFunc(fn)
}

func mapConfigHandler(cfg *mapConfig) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err := enc.Encode(cfg)

		if err != nil {
			slog.Error("Failed to encode map config", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	return http.HandlerFunc(fn)
}
