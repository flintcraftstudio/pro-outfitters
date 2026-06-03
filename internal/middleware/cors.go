package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig controls which cross-origin requests are permitted.
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int // seconds; only sent on preflight
}

func (c *CORSConfig) isOriginAllowed(origin string) bool {
	for _, allowed := range c.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

func (c *CORSConfig) allowedMethods() string {
	if len(c.AllowedMethods) == 0 {
		return "GET, POST, PUT, DELETE, OPTIONS"
	}
	return strings.Join(c.AllowedMethods, ", ")
}

func (c *CORSConfig) allowedHeaders() string {
	if len(c.AllowedHeaders) == 0 {
		return "Content-Type"
	}
	return strings.Join(c.AllowedHeaders, ", ")
}

// CORS returns middleware that handles cross-origin request headers.
// Disallowed origins receive no CORS headers — the browser enforces
// the block.
func CORS(config CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")
			if origin == "" || !config.isOriginAllowed(origin) {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				w.Header().Add("Vary", "Access-Control-Request-Method")
				w.Header().Add("Vary", "Access-Control-Request-Headers")
				w.Header().Set("Access-Control-Allow-Methods", config.allowedMethods())
				w.Header().Set("Access-Control-Allow-Headers", config.allowedHeaders())

				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
