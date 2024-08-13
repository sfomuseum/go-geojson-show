package show

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	MapProvider string
	MapTileURI  string
	Port        int
	URIs        []string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) *RunOptions {

	flagset.Parse(fs)

	uris := fs.Args()

	opts := &RunOptions{
		MapProvider: map_provider,
		MapTileURI:  map_tile_uri,
		Port:        port,
		URIs:        uris,
	}

	return opts
}
