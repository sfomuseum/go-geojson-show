package show

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/go-www-show/v2"
)

var port int

var browser_uri string

var map_provider string
var map_tile_uri string
var protomaps_theme string

var style string
var point_style string

var label_properties multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("show")

	browser_schemes := show.BrowserSchemes()
	str_schemes := strings.Join(browser_schemes, ",")

	browser_desc := fmt.Sprintf("A valid sfomuseum/go-www-show/v2.Browser URI. Valid options are: %s", str_schemes)

	fs.StringVar(&browser_uri, "browser-uri", "web://", browser_desc)

	fs.StringVar(&map_provider, "map-provider", "leaflet", "Valid options are: leaflet, protomaps")
	fs.StringVar(&map_tile_uri, "map-tile-uri", leaflet_osm_tile_url, "A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs.")
	fs.StringVar(&protomaps_theme, "protomaps-theme", "white", "A valid Protomaps theme label.")

	fs.StringVar(&style, "style", "", "A custom Leaflet style definition for geometries. This may either be a JSON-encoded string or a path on disk.")
	fs.StringVar(&point_style, "point-style", "", "A custom Leaflet style definition for point geometries. This may either be a JSON-encoded string or a path on disk.")
	fs.IntVar(&port, "port", 0, "The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.")

	fs.Var(&label_properties, "label", "Zero or more (GeoJSON Feature) properties to use to construct a label for a feature's popup menu when it is clicked on.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for serving GeoJSON features from an on-demand web server.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf the only path as input is \"-\" then data will be read from STDIN.\n\n")
	}

	return fs
}
