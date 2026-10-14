package view

import (
	"context"
	"time"
)

// SiteName is the display name used in templates. main.go sets this
// from config.AppName at startup.
var SiteName = "Standard Template"

// Tracking IDs and Turnstile site key, set once at startup from config.
var (
	PixelID          string
	GtagID           string
	TurnstileSiteKey string
)

// BaseURL is the site's canonical origin (e.g. https://example.com), set
// once at startup from config.AppBaseURL. Used to build absolute URLs for
// Open Graph / Twitter Card tags, which scrapers require to be absolute.
var BaseURL string

// Indexable mirrors config.AllowIndexing. When false every page carries
// noindex, so a staging deploy stays out of search results.
var Indexable bool

type pathKey struct{}

// WithPath stores the request path on ctx so the layout can emit a
// canonical URL without every handler passing it through.
func WithPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, pathKey{}, path)
}

// canonicalURL is BaseURL plus the request path stored by WithPath.
func canonicalURL(ctx context.Context) string {
	p, _ := ctx.Value(pathKey{}).(string)
	if p == "" {
		p = "/"
	}
	return BaseURL + p
}

// OGImage is the path (under BaseURL) to the default social share image.
const OGImage = "/static/img/hero-fishing-1200w.jpg"

// Year returns the current year for copyright notices.
func Year() int {
	return time.Now().Year()
}
