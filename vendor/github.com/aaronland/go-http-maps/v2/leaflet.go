package maps

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const LEAFLET_OSM_TILE_URL = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"

type LeafletConfig struct {
	Style      *LeafletStyle `json:"style,omitempty"`
	PointStyle *LeafletStyle `json:"point_style,omitempty"`
}

// LeafletStyle is a struct containing details for decorating GeoJSON features and markers
type LeafletStyle struct {
	Color       string  `json:"color,omitempty"`
	FillColor   string  `json:"fillColor,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
	Opacity     float64 `json:"opacity,omitempty"`
	Radius      float64 `json:"radius,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
}

// UnmarshalLeafletStyle derives a `LeafletStyle` instance from 'raw'. If 'raw' starts with "{" then it is treated as
// a JSON-encoded string, otherwise it is treated as a local path on disk.
func UnmarshalLeafletStyle(raw string) (*LeafletStyle, error) {

	raw = strings.TrimSpace(raw)

	if len(raw) == 0 {
		return nil, fmt.Errorf("Empty style definition")
	}

	if string(raw[0]) == "{" {
		return UnmarshalStyleFromString(raw)
	}

	r, err := os.Open(raw)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	return UnmarshalStyleFromReader(r)
}

// UnmarshalStyleFromString derives a `LeafletStyle` instance from 'raw'.
func UnmarshalStyleFromString(raw string) (*LeafletStyle, error) {

	var s *LeafletStyle

	err := json.Unmarshal([]byte(raw), &s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// UnmarshalStyleFromString derives a `LeafletStyle` instance from the body of 'r'.
func UnmarshalStyleFromReader(r io.Reader) (*LeafletStyle, error) {

	var s *LeafletStyle

	dec := json.NewDecoder(r)
	err := dec.Decode(&s)

	if err != nil {
		return nil, err
	}

	return s, nil
}
