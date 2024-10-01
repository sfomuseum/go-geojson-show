package show

import (
	"context"
	"flag"
	"fmt"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-flags/flagset"
	www_show "github.com/sfomuseum/go-www-show/v2"
)

type RunOptions struct {
	MapProvider     string
	MapTileURI      string
	ProtomapsTheme  string
	Port            int
	Features        []*geojson.Feature
	Style           *LeafletStyle
	PointStyle      *LeafletStyle
	LabelProperties []string
	Browser         www_show.Browser
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		MapProvider:     map_provider,
		MapTileURI:      map_tile_uri,
		ProtomapsTheme:  protomaps_theme,
		Port:            port,
		LabelProperties: label_properties,
	}

	br, err := www_show.NewBrowser(ctx, browser_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new browser, %w", err)
	}

	opts.Browser = br

	if style != "" {

		s, err := UnmarshalStyle(style)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal style, %w", err)
		}

		opts.Style = s
	}

	if point_style != "" {

		s, err := UnmarshalStyle(point_style)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal point style, %w", err)
		}

		opts.PointStyle = s
	}

	return opts, nil
}
