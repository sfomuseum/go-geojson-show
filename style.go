package show

type LeafletStyle struct {
	Color       string  `json:"color,omitempty"`
	FillColor   string  `json:"fillColor,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
	Opacity     float64 `json:"opacity,omitempty"`
	Radius      float64 `json:"radius,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
}
