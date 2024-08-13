// `go-http-leaflet` is an HTTP middleware package for including Leaflet.js assets in web applications. It exports two principal methods:
//
// * `leaflet.AppendAssetHandlers(*http.ServeMux)` which is used to append HTTP handlers to a `http.ServeMux` instance for serving Leaflet JavaScript and CSS files, and related assets.
// * `leaflet.AppendResourcesHandler(http.Handler, *LeafletOptions)` which is used to rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet JavaScript files and related assets.
//
// Example (Note the way we are embedding the HTML as an embed.FS instance)
//
//	import (
//		"embed"
//		"github.com/aaronland/go-http-leaflet"
//		"log"
//		"net/http"
//	)
//
//	//go:embed index.html
//	var staticFS embed.FS
//
//	func main() {
//
//		static_fs := http.FS(staticFS)
//		static_handler := http.FileServer(static_fs)
//
//		mux := http.NewServeMux()
//
//		leaflet.AppendAssetHandlers(mux)
//
//		leaflet_opts := leaflet.DefaultLeafletOptions()
//
//		index_handler := leaflet.AppendResourcesHandler(static_handler, leaflet_opts)
//		mux.Handle("/", index_handler)
//
//		err = http.ListenAndServe(":8080", mux)
//	}
package leaflet
