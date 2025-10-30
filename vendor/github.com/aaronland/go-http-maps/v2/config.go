package maps

import (
	"fmt"
)

type InitialView [2]float64

func (v *InitialView) String() string {
	return fmt.Sprintf("%f,%f", v[0], v[1])
}

type InitialBounds [4]float64

func (v *InitialBounds) String() string {
	return fmt.Sprintf("%f,%f,%f,%f", v[0], v[1], v[2], v[3])
}

// MapConfig defines common configuration details for maps.
type MapConfig struct {
	// A valid map provider label.
	Provider string `json:"provider"`
	// A valid Leaflet tile layer URI.
	TileURL string `json:"tile_url"`
	// Optional Protomaps configuration details
	Protomaps *ProtomapsConfig `json:"protomaps,omitempty"`
	// Optional Leaflet configuration details
	Leaflet *LeafletConfig `json:"leaflet,omitempty"`
	// The initial view (lon, lat) for the map (optional).
	InitialView *InitialView `json:"initial_view,omitempty"`
	// The initial zoom level for the map (optional).
	InitialZoom int `json:"initial_zoom,omitempty"`
	// The initial bounds (minx, miny, maxx, maxy) for the map (optional).
	InitialBounds *InitialBounds `json:"initial_bounds,omitempty"`
}
