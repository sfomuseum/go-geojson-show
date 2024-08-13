// `go-http-protomaps` is an HTTP middleware package for including Protomaps.js assets in web applications. It exports two principal methods:
//
// * `protomaps.AppendAssetHandlers(*http.ServeMux)` which is used to append HTTP handlers to a `http.ServeMux` instance for serving Protomaps JavaScript files, and related assets.
// * `protomaps.AppendResourcesHandler(http.Handler, *ProtomapsOptions)` which is used to rewrite any HTML produced by previous handler to include the necessary markup to load Protomaps JavaScript files and related assets.
//
// Example (Note the way we are embedding the HTML and .pmtiles database as an embed.FS instance)
//
//	import (
//		"embed"
//		"github.com/sfomuseum/go-http-protomaps"
//		"log"
//		"net/http"
//		"net/url"
//	)
//
//	//go:embed index.html sfo.pmtiles
//	var staticFS embed.FS
//
//	func main() {
//
//		tile_url := "/sfo.pmtiles"
//
//		static_fs := http.FS(staticFS)
//		static_handler := http.FileServer(static_fs)
//
//		mux := http.NewServeMux()
//
//		mux.Handle(*tile_url, static_handler)
//
//		protomaps.AppendAssetHandlers(mux)
//
//		pm_opts := protomaps.DefaultProtomapsOptions()
//		pm_opts.TileURL = *tile_url
//
//		index_handler := protomaps.AppendResourcesHandler(static_handler, pm_opts)
//		mux.Handle("/", index_handler)
//
//		err = http.ListenAndServe(":8080", mux)
//	}
package protomaps
