package show

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/aaronland/go-http-maps/v2"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-geojson-show/static/www"
	www_show "github.com/sfomuseum/go-www-show/v2"
	"github.com/tidwall/gjson"
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

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}
	mux := http.NewServeMux()

	www_fs := http.FS(www.FS)
	mux.Handle("/", http.FileServer(www_fs))

	fc := geojson.NewFeatureCollection()
	fc.Features = opts.Features
	data_handler := dataHandler(fc)

	mux.Handle("/features.geojson", data_handler)

	map_opts := &maps.AssignMapConfigHandlerOptions{
		MapProvider:            opts.MapProvider,
		MapTileURI:             opts.MapTileURI,
		LeafletStyle:           opts.Style,
		LeafletPointStyle:      opts.PointStyle,
		LeafletLabelProperties: opts.LabelProperties,
		LeafletPanes:           opts.LeafletPanes,
		ProtomapsTheme:         opts.ProtomapsTheme,
		ProtomapsMaxDataZoom:         opts.ProtomapsMaxDataZoom,		
	}

	err := maps.AssignMapConfigHandler(map_opts, mux, "/map.json")

	if err != nil {
		return fmt.Errorf("Failed to assign map config handler, %w", err)
	}

	local_cfg := &LocalConfig{
		ClusterMarkers: opts.ClusterMarkers,
	}

	config_handler := LocalConfigHandler(local_cfg)
	mux.Handle("/config.json", config_handler)

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
