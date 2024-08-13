package show

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var port int

var map_provider string
var map_tile_uri string
var protomaps_theme string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("show")

	fs.StringVar(&map_provider, "map-provider", "leaflet", "Valid options are: leaflet, protomaps")
	fs.StringVar(&map_tile_uri, "map-tile-uri", leaflet_osm_tile_url, "A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs.")
	fs.StringVar(&protomaps_theme, "protomaps-theme", "white", "A valid Protomaps theme label.")

	fs.IntVar(&port, "port", 0, "The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for serving GeoJSON files from an on-demand web server.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf the only path as input is \"-\" then data will be read from STDIN.\n\n")
	}

	return fs
}
