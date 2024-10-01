package show

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-geojson-show/static/www"
	"github.com/sfomuseum/go-http-protomaps"
	www_show "github.com/sfomuseum/go-www-show/v2"
	"github.com/tidwall/gjson"
	wasm_js "github.com/whosonfirst/go-whosonfirst-format-wasm/static/javascript"
	"github.com/whosonfirst/go-whosonfirst-format-wasm/static/wasm"
)

const leaflet_osm_tile_url = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
const protomaps_api_tile_url string = "https://api.protomaps.com/tiles/v3/{z}/{x}/{y}.mvt?key={key}"

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	fs_uris := fs.Args()

	features := make([]*geojson.Feature, 0)

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

			features = append(features, f)

		case "FeatureCollection":

			other_fc, err := geojson.UnmarshalFeatureCollection(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal record as FeatureCollection, %w", err)
			}

			for _, f := range other_fc.Features {
				features = append(features, f)
			}

		default:
			return fmt.Errorf("Invalid type, %s", type_rsp.String())
		}

		return nil
	}

	stdin := false

	if len(fs_uris) == 1 && fs_uris[0] == "-" {
		stdin = true
	}

	if stdin {

		err := append_features(os.Stdin)

		if err != nil {
			return fmt.Errorf("Failed to append features, %v", err)
		}

	} else {

		for _, path := range fs_uris {

			r, err := os.Open(path)

			if err != nil {
				return fmt.Errorf("Failed to open %s for reading, %v", path, err)
			}

			defer r.Close()

			err = append_features(r)

			if err != nil {
				return fmt.Errorf("Failed to append features, %v", err)
			}
		}
	}

	opts.Features = features

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	mux := http.NewServeMux()

	www_fs := http.FS(www.FS)
	mux.Handle("/", http.FileServer(www_fs))

	wasm_fs := http.FS(wasm.FS)
	wasm_handler := http.FileServer(wasm_fs)

	wasm_js_fs := http.FS(wasm_js.FS)
	wasm_js_handler := http.FileServer(wasm_js_fs)

	mux.Handle("/javascript/wasm/", http.StripPrefix("/javascript/wasm/", wasm_js_handler))
	mux.Handle("/wasm/", http.StripPrefix("/wasm/", wasm_handler))

	fc := geojson.NewFeatureCollection()
	fc.Features = opts.Features
	data_handler := dataHandler(fc)

	mux.Handle("/features.geojson", data_handler)

	//

	map_cfg := &mapConfig{
		Provider:        opts.MapProvider,
		TileURL:         opts.MapTileURI,
		Style:           opts.Style,
		PointStyle:      opts.PointStyle,
		LabelProperties: opts.LabelProperties,
	}

	if map_provider == "protomaps" {

		u, err := url.Parse(opts.MapTileURI)

		if err != nil {
			log.Fatalf("Failed to parse Protomaps tile URL, %w", err)
		}

		switch u.Scheme {
		case "file":

			mux_url, mux_handler, err := protomaps.FileHandlerFromPath(u.Path, "")

			if err != nil {
				log.Fatalf("Failed to determine absolute path for '%s', %v", opts.MapTileURI, err)
			}

			mux.Handle(mux_url, mux_handler)
			map_cfg.TileURL = mux_url

		case "api":
			key := u.Host
			map_cfg.TileURL = strings.Replace(protomaps_api_tile_url, "{key}", key, 1)
		}

		map_cfg.Protomaps = &protomapsConfig{
			Theme: opts.ProtomapsTheme,
		}
	}

	map_cfg_handler := mapConfigHandler(map_cfg)

	mux.Handle("/map.json", map_cfg_handler)

	www_show_opts := &www_show.RunOptions{
		Port:    opts.Port,
		Browser: opts.Browser,
		Mux:     mux,
	}

	return www_show.RunWithOptions(ctx, www_show_opts)
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
