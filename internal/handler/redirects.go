package handler

import (
	"net/http"
	"strings"
)

// legacyRedirects maps paths from the old WordPress site (taken from its
// Yoast sitemaps and homepage links, October 2026) to their new homes.
// Keys have no trailing slash; the old site used one on every URL and the
// lookup trims it. Old paths that match a new route exactly (/about) need
// no entry. Feeds, author, testimonial and slider archives are left to
// 404: nothing links to them and they have no counterpart.
var legacyRedirects = map[string]string{
	// Pages.
	"/lodges-rivers":                   "/lodges",
	"/contact-pro-outfitters":          "/contact",
	"/news-videos":                     "/",
	"/trip-details-nforkcrossinglodge": "/lodges/north-fork-crossing",
	"/smith-river-packing-list":        "/packing-lists/smith-river",
	"/kids-camp-what-to-bring":         "/packing-lists/kids-camp",
	"/bird-hunting-equipment-list":     "/packing-lists/bird-hunting",

	// Portfolio pages (lodges and rivers).
	"/portfolio_page":                           "/lodges",
	"/portfolio_page/smith-river-montana":       "/smith-river",
	"/portfolio_page/north-fork-crossing-lodge": "/lodges/north-fork-crossing",
	"/portfolio_page/sharptail-lodge":           "/lodges/sharptail",
	"/portfolio_page/yurt-craig":                "/lodges/yurt-at-craig",
	"/portfolio_page/kids-camp-prooutfitters":   "/kids-camp",
	"/portfolio_page/blackfoot-river":           "/blackfoot-river",

	// Blog posts. No news section on the new site; each goes to the page
	// it was about.
	"/2020/08/08/sharptail-lodge-wingshooting":            "/lodges/sharptail",
	"/2020/08/08/community-escape-smith-river-float-trip": "/smith-river",
	"/2020/07/16/sharptail-grouse-lodge":                  "/lodges/sharptail",
	"/2021/12/02/yeti-sharptail-taupe-collection":         "/lodges/sharptail",

	// Category archives.
	"/category/bird-hunting":        "/bird-hunting",
	"/category/trout-fishing":       "/fly-fishing",
	"/category/smith-river-montana": "/smith-river",
	"/category/sharptail-lodge":     "/lodges/sharptail",
	"/portfolio-category/featured":  "/",
	"/portfolio-category/lodges":    "/lodges",
	"/portfolio-category/rivers":    "/fly-fishing",
}

// LegacyRedirects answers old-site paths with a single 301 to the new
// page. It must wrap StripTrailingSlash so an old URL like
// /portfolio_page/sharptail-lodge/ doesn't take a 308 hop first. GET and
// HEAD only; the query string is dropped, since WordPress query params
// mean nothing here.
func LegacyRedirects(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			p := strings.TrimSuffix(strings.ToLower(r.URL.Path), "/")
			if to, ok := legacyRedirects[p]; ok {
				http.Redirect(w, r, to, http.StatusMovedPermanently)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
