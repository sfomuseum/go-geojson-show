package maps

import (
	"fmt"
	"net/http"
	"path/filepath"
)

const PROTOMAPS_API_TILE_URL = "https://api.protomaps.com/tiles/v3/{z}/{x}/{y}.mvt?key={key}"

// ProtomapsConfig defines configuration details for maps using Protomaps.
type ProtomapsConfig struct {
	// A valid Protomaps theme label
	Theme string `json:"theme"`
}

// ProtomapsFileHandlerFromPath will take a path and create a http.FileServer handler
// instance for the files in its root directory. The handler is returned with
// a relative URI for the filename in 'path' to be assigned to a net/http
// ServeMux instance.
func ProtomapsFileHandlerFromPath(path string, prefix string) (string, http.Handler, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return "", nil, fmt.Errorf("Failed to determine absolute path for '%s', %v", path, err)
	}

	fname := filepath.Base(abs_path)
	root := filepath.Dir(abs_path)

	tile_dir := http.Dir(root)
	tile_handler := http.FileServer(tile_dir)

	tile_url := fmt.Sprintf("/%s", fname)

	if prefix != "" {
		tile_handler = http.StripPrefix(prefix, tile_handler)
		tile_url = filepath.Join(prefix, fname)
	}

	return tile_url, tile_handler, nil
}
