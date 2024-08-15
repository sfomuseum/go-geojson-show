package show

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/tidwall/gjson"
)

type RunOptions struct {
	MapProvider    string
	MapTileURI     string
	ProtomapsTheme string
	Port           int
	Features       []*geojson.Feature
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		MapProvider:    map_provider,
		MapTileURI:     map_tile_uri,
		ProtomapsTheme: protomaps_theme,
		Port:           port,
	}

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

	uris := fs.Args()
	stdin := false

	if len(uris) == 1 && uris[0] == "-" {
		stdin = true
	}

	if stdin {

		err := append_features(os.Stdin)

		if err != nil {
			return nil, fmt.Errorf("Failed to append features, %v", err)
		}

	} else {

		for _, path := range uris {

			r, err := os.Open(path)

			if err != nil {
				return nil, fmt.Errorf("Failed to open %s for reading, %v", path, err)
			}

			defer r.Close()

			err = append_features(r)

			if err != nil {
				return nil, fmt.Errorf("Failed to append features, %v", err)
			}
		}
	}

	opts.Features = features
	return opts, nil
}
