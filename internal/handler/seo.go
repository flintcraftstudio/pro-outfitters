package handler

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flintcraftstudio/standard-template/internal/view"
)

// Robots handles GET /robots.txt. Until indexing is allowed it disallows
// everything, so a staging deploy never reaches search results.
func Robots(baseURL string, allowIndexing bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if !allowIndexing {
			fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
			return
		}
		fmt.Fprintf(w, "User-agent: *\nDisallow: /admin\nDisallow: /login\n\nSitemap: %s/sitemap.xml\n", baseURL)
	}
}

type urlset struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []sitemapURL
}

type sitemapURL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

// Sitemap handles GET /sitemap.xml. It lists the finished public pages
// only; draft pages (view.DraftBase) are noindex and stay out until their
// content arrives and they move to view.Base.
func Sitemap(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paths := []string{"/", "/lodges"}
		for _, l := range view.Lodges() {
			paths = append(paths, "/lodges/"+l.Slug)
		}
		paths = append(paths, "/smith-river", "/about", "/contact")

		set := urlset{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
		for _, p := range paths {
			set.URLs = append(set.URLs, sitemapURL{Loc: baseURL + p})
		}
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		fmt.Fprint(w, xml.Header)
		if err := xml.NewEncoder(w).Encode(set); err != nil {
			slog.Error("sitemap encode", "err", err)
		}
	}
}
