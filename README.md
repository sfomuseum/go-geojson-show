# go-geojson-show

Command-line tool for server GeoJSON files from an on-demand web server.

## Documentation

Documentation is incomplete at this time.

## Tools

### show

```
$> ./bin/show -h
Usage of ./bin/show:
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -port int
    	The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
```

#### Example

```
$> cat feature.geojson | ./bin/show  -
2024/08/12 18:00:26 Features are viewable at http://localhost:54420
```

```
$> cat feature.geojson | ./bin/show -map-provider leaflet -map-tile-url '...' -
2024/08/12 18:16:51 Features are viewable at http://localhost:61222
```

```
$> cat feature.geojson | ./bin/show -map-provider protomaps -map-tile-url file:///path/to/your/database.pmtiles -
2024/08/12 18:17:51 Features are viewable at http://localhost:49316
```

```
$> cat feature.geojson | ./bin/show -map-provider protomaps -map-tile-url api://{APIKEY} -
2024/08/12 18:18:14 Features are viewable at http://localhost:51021
```