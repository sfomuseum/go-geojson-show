package www

import (
	"embed"
)

//go:embed javascript/* css/* wasm/*.wasm *.html
var FS embed.FS
