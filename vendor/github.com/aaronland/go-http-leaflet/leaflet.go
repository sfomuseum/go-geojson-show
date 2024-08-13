package leaflet

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aaronland/go-http-leaflet/static"
	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-rollup"
)

var ROLLUP_ASSETS = false

// LeafletOptions provides a list of JavaScript and CSS link to include with HTML output.
type LeafletOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
	RollupAssets          bool
	Prefix                string
}

// Append the Javascript and CSS URLs for the Leaflet.Fullscreen plugin.
func (opts *LeafletOptions) EnableFullscreen() {
	opts.CSS = append(opts.CSS, "/css/leaflet.fullscreen.css")
	opts.JS = append(opts.JS, "/javascript/leaflet.fullscreen.min.js")
}

// Append the Javascript and CSS URLs for the Leaflet.Hash plugin.
func (opts *LeafletOptions) EnableHash() {
	opts.JS = append(opts.JS, "/javascript/leaflet-hash.js")
}

// Append the Javascript and CSS URLs for the leaflet-geoman plugin.
// https://github.com/geoman-io/leaflet-geoman/
func (opts *LeafletOptions) EnableDraw() {
	opts.CSS = append(opts.CSS, "/css/leaflet-geoman.css")
	opts.JS = append(opts.JS, "/javascript/leaflet-geoman.min.js")
}

// Return a *LeafletOptions struct with default paths and URIs.
func DefaultLeafletOptions() *LeafletOptions {

	opts := &LeafletOptions{
		CSS: []string{
			"/css/leaflet.css",
		},
		JS: []string{
			"/javascript/leaflet.js",
		},
		DataAttributes: make(map[string]string),
	}

	return opts
}

// AppendResourcesHandlerWithPrefix will rewrite any HTML produced by previous handler to include the necessary markup to load Leaflet JavaScript files and related assets ensuring that all URIs are prepended with a prefix.
func AppendResourcesHandler(next http.Handler, opts *LeafletOptions) http.Handler {

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.DataAttributes = opts.DataAttributes
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	if opts.RollupAssets {

		static_opts.CSS = []string{
			"/css/leaflet.rollup.css",
		}

		static_opts.JS = []string{
			"/javascript/leaflet.rollup.js",
		}

	} else {

		static_opts.CSS = opts.CSS
		static_opts.JS = opts.JS
	}

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, opts.Prefix)
}

// Append all the files in the net/http FS instance containing the embedded Leaflet assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux, opts *LeafletOptions) error {

	if !opts.RollupAssets {
		return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, opts.Prefix)
	}

	js_paths := make([]string, len(opts.JS))
	css_paths := make([]string, len(opts.CSS))

	for idx, path := range opts.JS {
		path = strings.TrimLeft(path, "/")
		js_paths[idx] = path
	}

	for idx, path := range opts.CSS {
		path = strings.TrimLeft(path, "/")
		css_paths[idx] = path
	}

	rollup_js_paths := map[string][]string{
		"leaflet.rollup.js": js_paths,
	}

	rollup_js_opts := &rollup.RollupJSHandlerOptions{
		FS:    static.FS,
		Paths: rollup_js_paths,
	}

	rollup_js_handler, err := rollup.RollupJSHandler(rollup_js_opts)

	if err != nil {
		return fmt.Errorf("Failed to create rollup JS handler, %w", err)
	}

	rollup_js_uri := "/javascript/leaflet.rollup.js"

	if opts.Prefix != "" {

		u, err := url.JoinPath(opts.Prefix, rollup_js_uri)

		if err != nil {
			return fmt.Errorf("Failed to append prefix to %s, %w", rollup_js_uri, err)
		}

		rollup_js_uri = u
	}

	mux.Handle(rollup_js_uri, rollup_js_handler)

	// CSS

	rollup_css_paths := map[string][]string{
		"leaflet.rollup.css": css_paths,
	}

	rollup_css_opts := &rollup.RollupCSSHandlerOptions{
		FS:    static.FS,
		Paths: rollup_css_paths,
	}

	rollup_css_handler, err := rollup.RollupCSSHandler(rollup_css_opts)

	if err != nil {
		return fmt.Errorf("Failed to create rollup CSS handler, %w", err)
	}

	rollup_css_uri := "/css/leaflet.rollup.css"

	if opts.Prefix != "" {

		u, err := url.JoinPath(opts.Prefix, rollup_css_uri)

		if err != nil {
			return fmt.Errorf("Failed to append prefix to %s, %w", rollup_css_uri, err)
		}

		rollup_css_uri = u
	}

	mux.Handle(rollup_css_uri, rollup_css_handler)
	return nil
}
