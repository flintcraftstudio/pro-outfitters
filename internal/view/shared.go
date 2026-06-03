package view

import "time"

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

// OGImage is the path (under BaseURL) to the default social share image.
const OGImage = "/static/img/hero-fishing-1200w.jpg"

// Year returns the current year for copyright notices.
func Year() int {
	return time.Now().Year()
}
