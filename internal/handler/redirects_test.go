package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLegacyRedirects(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	h := LegacyRedirects(next)

	tests := []struct {
		method, path string
		code         int
		location     string
	}{
		{"GET", "/portfolio_page/sharptail-lodge/", 301, "/lodges/sharptail"},
		{"GET", "/portfolio_page/sharptail-lodge", 301, "/lodges/sharptail"},
		{"HEAD", "/smith-river-packing-list/", 301, "/packing-lists/smith-river"},
		{"GET", "/Contact-Pro-Outfitters/?utm_source=x", 301, "/contact"},
		{"POST", "/contact-pro-outfitters/", http.StatusTeapot, ""},
		{"GET", "/about/", http.StatusTeapot, ""},
		{"GET", "/lodges", http.StatusTeapot, ""},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
		if rec.Code != tt.code {
			t.Errorf("%s %s: code = %d, want %d", tt.method, tt.path, rec.Code, tt.code)
		}
		if got := rec.Header().Get("Location"); got != tt.location {
			t.Errorf("%s %s: Location = %q, want %q", tt.method, tt.path, got, tt.location)
		}
	}
}

// Every redirect target must be a real route, or the redirect lands on a
// 404. Keep this list in step with cmd/server/main.go.
func TestLegacyRedirectTargetsExist(t *testing.T) {
	routes := map[string]bool{
		"/": true, "/lodges": true, "/smith-river": true, "/about": true,
		"/contact": true, "/fly-fishing": true, "/bird-hunting": true,
		"/kids-camp": true, "/blackfoot-river": true,
		"/lodges/north-fork-crossing": true, "/lodges/sharptail": true,
		"/lodges/yurt-at-craig":      true,
		"/packing-lists/smith-river": true, "/packing-lists/kids-camp": true,
		"/packing-lists/bird-hunting": true,
	}
	for from, to := range legacyRedirects {
		if !routes[to] {
			t.Errorf("%s redirects to %s, which is not a known route", from, to)
		}
	}
}
