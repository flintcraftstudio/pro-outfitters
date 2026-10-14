package middleware

import "net/http"

// RealIP sets r.RemoteAddr to the client's IP, read from X-Forwarded-For
// past trustedProxies hops (see clientIP). Behind Caddy, RemoteAddr is
// otherwise the proxy's address, so request logs, Turnstile's remoteip
// and anything else that reads it would all see Caddy instead of the
// visitor. Wrap it outermost so every layer sees the rewritten value.
//
// Only trust the header when the app is reachable solely through the
// proxy (port bound to loopback); otherwise a client can forge it.
func RealIP(trustedProxies int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.RemoteAddr = clientIP(r, trustedProxies)
			next.ServeHTTP(w, r)
		})
	}
}
