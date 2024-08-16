# go-geojson-show

Command-line tool for serving GeoJSON features from an on-demand web server.

## Motivation

It's basically a simpler and dumber version of [geojson.io](https://geojson.io/) that you can run locally from a single binary application. Also, the option for custom, local and private tile data.

## Documentation

Documentation (`godoc`) is incomplete at this time.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/show cmd/show/main.go
```

### show

```
$> ./bin/show -h
Command-line tool for serving GeoJSON features from an on-demand web server.
Usage:
	 ./bin/show path(N) path(N)
Valid options are:
  -label value
    	Zero or more (GeoJSON Feature) properties to use to construct a label for a feature's popup menu when it is clicked on.
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -point-style string
    	A custom Leaflet style definition for point geometries. This may either be a JSON-encoded string or a path on disk.
  -port int
    	The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
  -style string
    	A custom Leaflet style definition for geometries. This may either be a JSON-encoded string or a path on disk.

If the only path as input is "-" then data will be read from STDIN.
```

#### Examples

##### Read a single GeoJSON file from disk and show it on a map using the default settings (OpenStreetMap)

![](docs/images/go-geojson-show-simple.png)

```
$> ./bin/show \
	/usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson
	
2024/08/13 13:01:44 Features are viewable at http://localhost:55799
```

##### Read multiple GeoJSON files from disk and show them on a map using the default settings (OpenStreetMap)

![](docs/images/go-geojson-show-multi.png)


```
$> ./bin/show \
	/usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson \
	/usr/local/data/oak.geojson
	
2024/08/13 13:08:44 Features are viewable at http://localhost:54501
```

##### Read a single GeoJSON file from disk and show it on a map using custom tiles:

![](docs/images/go-geojson-show-custom.png)

```
$> ./bin/show \
	-map-tile-uri 'https://static.sfomuseum.org/aerial/1978/{z}/{x}/{-y}.png'
	/usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson
	
2024/08/13 13:03:17 Features are viewable at http://localhost:62669
```

##### Read the (GeoJSON) output of another process and show those features on a map using a local [Protomaps](https://protomaps.com) database file and a named Protomaps theme

![](docs/images/go-geojson-show-protomaps-local.png)

```
$> cat /usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson | \
	./bin/show \
	-map-provider protomaps \
	-map-tile-uri file:///usr/local/sfomuseum/go-http-protomaps/cmd/example/sfo.pmtiles \
	-protomaps-theme light \
	-
	
2024/08/13 13:05:13 Features are viewable at http://localhost:54749
```

##### Read the (GeoJSON) output of another process and show those features on a map using the Protomaps API

![](docs/images/go-geojson-show-protomaps-api.png)

```
$> cat /usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson | \
	./bin/show \
	-map-provider protomaps \
	-map-tile-uri api://{APIKEY} \
	-
	
2024/08/13 13:07:08 Features are viewable at http://localhost:63818
```

##### Read a single GeoJSON file from disk and show it with a custom marker style

![](docs/images/go-geojson-show-styles.png)

```
$> ./bin/show \
	-point-style '{"radius": 10, "color": "red", "fillColor": "orange" }' \
	/usr/local/data/postcards.geojson
	
2024/08/15 16:15:37 Features are viewable at http://localhost:63516
```

See [styles.go](styles.go) for details about the structure of the `LeafletStyle` struct used to encode custom map styles.

##### Read a single GeoJSON file from disk and show it with custom labels when a marker is clicked

![](docs/images/go-geojson-show-label.png)

```
$> ./bin/show \
	-label wof:name \
	-label wof:id \
	/usr/local/data/postcards.geojson
	
2024/08/15 16:12:39 Features are viewable at http://localhost:50310
```

When a marker is clicked the application will scroll that feature's string representation (in the right-hand pane) in to view and highlight its text.

## Advanced usage

### Using `go-geojson-show` as a package

What follows is an annotated and abbreviated version of the code used by the [whosonfirst/wof-cli](https://github.com/whosonfirst/wof-cli/blob/main/show/show.go) package to show features on a map using the `sfomuseum/go-geojson-show` package.

_For the sake of brevity error handling has been omitted._

#### Step 1: Parsing flags and deriving default "run options":

The first step is to import any necessary packages including `github.com/sfomuseum/go-geojson-show` which is used to define a default flag set, parse command line arguments and then derive "run options" for the application.

```
import (
        "context"
	"io"
	"slices"

	"github.com/paulmach/orb/geojson"
	sfom_show "github.com/sfomuseum/go-geojson-show"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"	
)

func show(args []string) {

	fs := sfom_show.DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	run_opts, _ := sfom_show.RunOptionsFromFlagSet(ctx, fs)
```

#### Step 2: Doing custom work to derive a list of `geojson.Feature` records to display

This is custom code, specific to the `wof-cli` package. It defines a set of default properties to use for marker labels and supplements them with any new labels passed defined in the flagset / run options. Afterwards it derives one or more GeoJSON feature records, using its own internal logic, from paths defined on the command line.

```
	label_props := []string{
		"wof:name",
		"wof:id",
		"wof:placetype",
		"src:geom",
	}

	for _, prop := range run_opts.LabelProperties {

		if !slices.Contains(label_props, prop) {
			label_props = append(label_props, prop)
		}
	}

	run_opts.LabelProperties = label_props

	fc := geojson.NewFeatureCollection()

	cb := func(ctx context.Context, uri string) error {

		r, is_stdin, _ := reader.ReadCloserFromURI(ctx, uri)

		if !is_stdin {
			defer r.Close()
		}

		body, _ := io.ReadAll(r)

		f, _ := geojson.UnmarshalFeature(body)

		fc.Append(f)
		return nil
	}

	uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
```

#### Step 3: Showing features on a map

Finally, the run options are updated with the new list of features and the `RunWithOptions` method is invoked.

```
	run_opts.Features = fc.Features
	return sfom_show.RunWithOptions(ctx, run_opts)
}
```

That's it.