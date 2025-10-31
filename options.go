package show

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-flags/flagset"
	www_show "github.com/sfomuseum/go-www-show/v2"
)

type RunOptions struct {
	MapProvider     string
	MapTileURI      string
	ProtomapsTheme  string
	ProtomapsMaxDataZoom int
	Port            int
	Features        []*geojson.Feature
	Style           string
	PointStyle      string
	LabelProperties []string
	LeafletPanes    map[string]int
	ClusterMarkers  bool
	Browser         www_show.Browser
	Verbose         bool
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		MapProvider:     map_provider,
		MapTileURI:      map_tile_uri,
		ProtomapsTheme:  protomaps_theme,
		ProtomapsMaxDataZoom: protomaps_max_data_zoom,
		Port:            port,
		LabelProperties: label_properties,
		ClusterMarkers:  cluster_markers,
		Verbose:         verbose,
	}

	if len(panes) > 0 {

		leaflet_panes := make(map[string]int)

		for _, fl := range panes {
			leaflet_panes[fl.Key()] = int(fl.Value().(int64))
		}

		opts.LeafletPanes = leaflet_panes
	}

	br, err := www_show.NewBrowser(ctx, browser_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new browser, %w", err)
	}

	opts.Browser = br

	if style != "" {
		
		if !strings.HasPrefix(style, "{") {

			body, err := os.ReadFile(style)

			if err != nil {
				return nil, fmt.Errorf("Failed to read style definition, %w", err)
			}

			style = string(body)
		}

		opts.Style = style
	}

	if point_style != "" {
		
		if !strings.HasPrefix(point_style, "{") {

			body, err := os.ReadFile(point_style)

			if err != nil {
				return nil, fmt.Errorf("Failed to read point style definition, %w", err)
			}

			point_style = string(body)
		}

		opts.PointStyle = point_style
	}

	return opts, nil
}
