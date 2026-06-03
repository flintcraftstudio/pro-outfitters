package middleware

import "net/http"

// StripTrailingSlash redirects any request whose path ends in "/" (other
// than bare "/") to the same path without the trailing slash, using 308
// to preserve method and body. Go 1.22+ ServeMux treats `GET /foo` and
// `GET /foo/` as distinct patterns and doesn't auto-canonicalize.
func StripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if len(p) > 1 && p[len(p)-1] == '/' {
			trimmed := p[:len(p)-1]
			target := trimmed
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusPermanentRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
