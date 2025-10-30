package maps

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// MapConfigHandler returns a new `http.Handler` that will return a JSON-encoded version of 'cfg'.
func MapConfigHandler(cfg *MapConfig) http.Handler {

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

// AssignMapConfigHandlerOptions defines configuration options for the AssignMapConfigHandler method.
type AssignMapConfigHandlerOptions struct {
	// A string label indicating the map provider to use. Valid options are: leaflet, protomaps.
	MapProvider string
	// A valid Leaflet tile layer URI.
	MapTileURI string
	// A comma-separated string indicating the map's initial view. Valid options are: 'LON,LAT', 'LON,LAT,ZOOM' or 'MINX,MINY,MAXX,MAXY'.
	InitialView string
	// A custom Leaflet style definition for geometries. This may either be a JSON-encoded string or a path on disk.
	LeafletStyle string
	// A custom Leaflet style definition for points. This may either be a JSON-encoded string or a path on disk.
	LeafletPointStyle string
	// A valid Protomaps theme label.
	ProtomapsTheme string
}

// AssignMapConfigHandler derives a new `MapConfig` and corresponding `http.Handler` instance from 'opts' assigning
// to 'mux' at 'map_cfg_uri'.
func AssignMapConfigHandler(opts *AssignMapConfigHandlerOptions, mux *http.ServeMux, map_cfg_uri string) error {

	map_cfg := &MapConfig{
		Provider: opts.MapProvider,
		TileURL:  opts.MapTileURI,
	}

	if opts.InitialView != "" {

		parts := strings.Split(opts.InitialView, ",")

		switch len(parts) {
		case 2:

			lon, err := strconv.ParseFloat(parts[0], 10)

			if err != nil {
				return fmt.Errorf("Failed to parse longitude, %w", err)
			}

			lat, err := strconv.ParseFloat(parts[1], 10)

			if err != nil {
				return fmt.Errorf("Failed to parse latitude, %w", err)
			}

			map_cfg.InitialView = &InitialView{lon, lat}

		case 3:

			lon, err := strconv.ParseFloat(parts[0], 10)

			if err != nil {
				return fmt.Errorf("Failed to parse longitude, %w", err)
			}

			lat, err := strconv.ParseFloat(parts[1], 10)

			if err != nil {
				return fmt.Errorf("Failed to parse latitude, %w", err)
			}

			zoom, err := strconv.Atoi(parts[2])

			if err != nil {
				return fmt.Errorf("Failed to parse zoom, %w", err)
			}

			map_cfg.InitialView = &InitialView{lon, lat}
			map_cfg.InitialZoom = zoom

		case 4:

			minx, err := strconv.ParseFloat(parts[0], 10)

			if err != nil {
				return fmt.Errorf("Invalid minx, %w", err)
			}

			miny, err := strconv.ParseFloat(parts[1], 10)

			if err != nil {
				return fmt.Errorf("Invalid miny, %w", err)
			}

			maxx, err := strconv.ParseFloat(parts[2], 10)

			if err != nil {
				return fmt.Errorf("Invalid maxx, %w", err)
			}

			maxy, err := strconv.ParseFloat(parts[3], 10)

			if err != nil {
				return fmt.Errorf("Invalid maxy, %w", err)
			}

			map_cfg.InitialBounds = &InitialBounds{minx, miny, maxx, maxy}

		default:
			return fmt.Errorf("Invalid initial view. Must be: 'lon,lat' or 'lon,lat,zoom' or 'minx,miny,maxx, maxy'")
		}
	}

	switch opts.MapProvider {
	case "leaflet":

		if opts.LeafletStyle != "" && opts.LeafletPointStyle != "" {

			leaflet_cfg := &LeafletConfig{}

			if opts.LeafletStyle != "" {

				s, err := UnmarshalLeafletStyle(opts.LeafletStyle)

				if err != nil {
					return fmt.Errorf("Failed to unmarshal leaflet style, %w", err)
				}

				leaflet_cfg.Style = s
			}

			if opts.LeafletPointStyle != "" {

				s, err := UnmarshalLeafletStyle(opts.LeafletPointStyle)

				if err != nil {
					return fmt.Errorf("Failed to unmarshal leaflet point style, %w", err)
				}

				leaflet_cfg.PointStyle = s

			}

			map_cfg.Leaflet = leaflet_cfg
		}

	case "protomaps":

		u, err := url.Parse(opts.MapTileURI)

		if err != nil {
			return fmt.Errorf("Failed to parse Protomaps tile URL, %w", err)
		}

		switch u.Scheme {
		case "file":

			mux_url, mux_handler, err := ProtomapsFileHandlerFromPath(u.Path, "")

			if err != nil {
				return fmt.Errorf("Failed to determine absolute path for '%s', %w", opts.MapTileURI, err)
			}

			mux.Handle(mux_url, mux_handler)
			map_cfg.TileURL = mux_url

		case "api":
			key := u.Host
			map_cfg.TileURL = strings.Replace(PROTOMAPS_API_TILE_URL, "{key}", key, 1)
		}

		map_cfg.Protomaps = &ProtomapsConfig{
			Theme: opts.ProtomapsTheme,
		}
	}

	map_cfg_handler := MapConfigHandler(map_cfg)
	mux.Handle(map_cfg_uri, map_cfg_handler)

	return nil
}
