package middleware

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitConfig controls per-client request throttling.
type RateLimitConfig struct {
	Rate            rate.Limit
	Burst           int
	TrustedProxies  int
	KeyFunc         func(r *http.Request) string
	CleanupInterval time.Duration
}

// RateLimit returns middleware that throttles requests per client
// using a token bucket algorithm.
func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	if config.Burst == 0 {
		config.Burst = 1
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 5 * time.Minute
	}

	store := newVisitorStore(config.Rate, config.Burst)
	go store.cleanup(config.CleanupInterval)

	keyFunc := config.KeyFunc
	if keyFunc == nil {
		keyFunc = func(r *http.Request) string {
			return clientIP(r, config.TrustedProxies)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			limiter := store.get(key)

			w.Header().Set("X-RateLimit-Limit", formatRate(config.Rate))
			w.Header().Set("X-RateLimit-Burst", strconv.Itoa(config.Burst))

			if !limiter.Allow() {
				retryAfter := int(math.Ceil(1.0 / float64(config.Rate)))
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the real client IP, respecting the number of
// trusted proxies. With TrustedProxies == 1 (e.g. one Caddy in front),
// take the rightmost X-Forwarded-For entry that the proxy appended.
// With TrustedProxies == 0 use RemoteAddr directly, ignoring the header.
func clientIP(r *http.Request, trustedProxies int) string {
	if trustedProxies > 0 {
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			parts := strings.Split(xff, ",")
			idx := max(len(parts)-trustedProxies, 0)
			ip := strings.TrimSpace(parts[idx])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type visitorStore struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	burst    int
}

func newVisitorStore(r rate.Limit, burst int) *visitorStore {
	return &visitorStore{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
	}
}

func (s *visitorStore) get(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.visitors[key]
	if !ok {
		v = &visitor{
			limiter: rate.NewLimiter(s.rate, s.burst),
		}
		s.visitors[key] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (s *visitorStore) cleanup(interval time.Duration) {
	for {
		time.Sleep(interval)
		cutoff := time.Now().Add(-2 * interval)
		s.mu.Lock()
		for key, v := range s.visitors {
			if v.lastSeen.Before(cutoff) {
				delete(s.visitors, key)
			}
		}
		s.mu.Unlock()
	}
}

func formatRate(r rate.Limit) string {
	if r == rate.Inf {
		return "unlimited"
	}
	f := float64(r)
	if f == math.Trunc(f) {
		return strconv.Itoa(int(f))
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}
