# go-http-maps

Go package providing opinionated HTTP middleware for web-based map tiles.

## Version 2

`/v2` of this package is not just a complete refactoring of the code but also a complete rethink about what its function is and how it works. There is basically nothing which is backwards compatible with previous versions.

## Motivation

The original motivation for the `go-http-maps` package was to define a handful of top-level methods and `net/http` middleware handlers to manage the drudegry of setting up maps, with a variety of map providers (Leaflet, Protomaps, Nextzen/Tangramjs), It "worked" but, in the end, it was also not "easy".

Version 2 removes most of the functionality of the `go-http-maps` package and instead focuses on a handful of methods for providing a dynamic map "config" file (exposed as an HTTP endpoint) which can be from the browser.

This package no longer provides static asset handlers for Leaflet or Protomaps. It is left up to you to bundle (and serve) them from your own code. You could, if you wanted, define a custom `http.Handler` instance to load those files from the `github.com/aaronland/go-http-maps/v2/static/www.FS` embedded filesystem but that's still something you'll need to do on your own.

## Documentation

`godoc` is still incomplete at this time.

## Usage

The easiest usage example is the [cmd/example](cmd/example/main.go) tool:

```
package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-maps/v2"
	"github.com/aaronland/go-http-maps/v2/static/www"
)

func main() {

	var verbose bool
	var host string
	var port int

	var initial_view string
	var map_provider string
	var map_tile_uri string
	var protomaps_theme string
	var leaflet_style string
	var leaflet_point_style string

	flag.StringVar(&host, "host", "localhost", "The host to listen for requests on")
	flag.IntVar(&port, "port", 8080, "The port number to listen for requests on")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	flag.StringVar(&map_provider, "map-provider", "leaflet", "Valid options are: leaflet, protomaps")
	flag.StringVar(&map_tile_uri, "map-tile-uri", maps.LEAFLET_OSM_TILE_URL, "A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs.")
	flag.StringVar(&protomaps_theme, "protomaps-theme", "white", "A valid Protomaps theme label.")
	flag.StringVar(&leaflet_style, "leaflet_style", "", "A custom Leaflet style definition for geometries. This may either be a JSON-encoded string or a path on disk.")
	flag.StringVar(&leaflet_point_style, "leaflet_point_style", "", "A custom Leaflet style definition for points. This may either be a JSON-encoded string or a path on disk.")
	flag.StringVar(&initial_view, "initial-view", "", "A comma-separated string indicating the map's initial view. Valid options are: 'LON,LAT', 'LON,LAT,ZOOM' or 'MINX,MINY,MAXX,MAXY'.")

	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	mux := http.NewServeMux()

	opts := &maps.AssignMapConfigHandlerOptions{
		MapProvider:       map_provider,
		MapTileURI:        map_tile_uri,
		InitialView:       initial_view,
		LeafletStyle:      leaflet_style,
		LeafletPointStyle: leaflet_point_style,
		ProtomapsTheme:    protomaps_theme,
	}

	maps.AssignMapConfigHandler(opts, mux, "/map.json")
	
	www_fs := http.FS(www.FS)
	www_handler := http.FileServer(www_fs)

	mux.Handle("/", www_handler)

	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("Listening for requests", "address", addr)

	http.ListenAndServe(addr, mux)
}
```

_Error handling omitted for the sake of brevity._

The "nut" of it being this part:

```
	opts := &maps.AssignMapConfigHandlerOptions{
		MapProvider:       map_provider,
		MapTileURI:        map_tile_uri,
		InitialView:       initial_view,
		LeafletStyle:      leaflet_style,
		LeafletPointStyle: leaflet_point_style,
		ProtomapsTheme:    protomaps_theme,
	}

	maps.AssignMapConfigHandler(opts, mux, "/map.json")
```

Which populates the `maps.AssignMapConfigHandlerOptions` with map-specific command line flags and passes those options (along with your `http.ServeMux` instance) to the `AssignMapConfigHandler` ethod. This method will validate all the options and create a new map config `http.Handler` assigning it to the `http.ServeMux` instance.

If you are using Protomaps as your map provider and the corresponding `MapTileURI` starts with `file://` it will be assumed that you are trying to serve a local Protomaps database and a matching `http.Handler` will assign to the `http.ServeMux` instance.

And then in your JavaScript code you would write something like this:

```
window.addEventListener("load", function load(event){

    // Null Island    
    var map = L.map('map').setView([0.0, 0.0], 1);    

    fetch("/map.json")
        .then((rsp) => rsp.json())
        .then((cfg) => {

	    console.debug("Got map config", cfg);
	    
            switch (cfg.provider) {
                case "leaflet":

                    var tile_url = cfg.tile_url;

                    var tile_layer = L.tileLayer(tile_url, {
                        maxZoom: 19,
                    });

                    tile_layer.addTo(map);
                    break;

                case "protomaps":

                    var tile_url = cfg.tile_url;

                    var tile_layer = protomapsL.leafletLayer({
                        url: tile_url,
                        theme: cfg.protomaps.theme,
                    })

                    tile_layer.addTo(map);
                    break;

                default:
                    console.error("Uknown or unsupported map provider");
                    return;
	    }

	    if (cfg.initial_view) {

		var zm = map.getZoom();

		if (cfg.initial_zoom){
		    zm = cfg.initial_zoom;
		}

		map.setView([cfg.initial_view[1], cfg.initial_view[0]], zm);
		
	    } else if (cfg.initial_bounds){

		var bounds = [
		    [ cfg.initial_bounds[1], cfg.initial_bounds[0] ],
		    [ cfg.initial_bounds[3], cfg.initial_bounds[2] ],
		];

		map.fitBounds(bounds);
	    }
	    
	    console.debug("Finished map setup");
	    
        }).catch((err) => {
	    console.error("Failed to derive map config", err);
	    return;
	});    
    
});
```

## Example

### Leaflet

![](docs/images/go-http-maps-leaflet.png)

```
$> make example
go run cmd/example/main.go \
		-initial-view '-122.384292,37.621131,13'
		
2025/03/06 10:00:09 INFO Listening for requests address=localhost:8080
```

### Protomaps

![](docs/images/go-http-maps-protomaps.png)

```
$> make example-protomaps
go run cmd/example/main.go \
		-initial-view '-122.384292,37.621131,13' \
		-map-provider protomaps \
		-map-tile-uri 'file:///usr/local/go-http-maps/fixtures/sfo.pmtiles'
		
2025/03/06 09:59:05 INFO Listening for requests address=localhost:8080
```

