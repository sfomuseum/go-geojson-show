package www

import (
	"embed"
)

//go:embed javascript/* css/* *.html
var FS embed.FS
